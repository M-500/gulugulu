package comment

import (
	"context"
	"errors"
	"fmt"
	"path/filepath"
	"strings"
	"time"

	commentrepo "gl-app/api/internal/repo/comment_repo"
	mediarepo "gl-app/api/internal/repo/media_repo"
	"gl-app/api/internal/svc"
	"gl-app/api/internal/types"
	"gl-app/pkg/ctxdata"

	"github.com/google/uuid"
	"github.com/minio/minio-go/v7"
	"github.com/zeromicro/go-zero/core/logx"
)

type CreateCommentLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCreateCommentLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateCommentLogic {
	return &CreateCommentLogic{Logger: logx.WithContext(ctx), ctx: ctx, svcCtx: svcCtx}
}

func (l *CreateCommentLogic) CreateComment(req *types.CreateCommentReq) (*types.CreateCommentResp, error) {
	userID := ctxdata.GetUidFromCtx(l.ctx)
	if userID <= 0 {
		return nil, fmt.Errorf("请先登录")
	}
	if req == nil || req.WorkId <= 0 {
		return nil, fmt.Errorf("作品ID不正确")
	}
	req.Content = strings.TrimSpace(req.Content)
	if req.Content == "" && req.ImageMediaId <= 0 {
		return nil, fmt.Errorf("评论文字和图片不能同时为空")
	}
	if len([]rune(req.Content)) > 500 {
		return nil, fmt.Errorf("评论内容最多500个字符")
	}
	if req.ParentCommentId < 0 {
		return nil, fmt.Errorf("回复评论ID不正确")
	}
	user, err := l.svcCtx.UserRepo.FindOneByID(l.ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("评论用户不存在")
	}

	input := commentrepo.CreateInput{WorkID: req.WorkId, UserID: userID, ParentID: req.ParentCommentId, Content: req.Content, ImageMediaID: req.ImageMediaId}
	var sourceBucket, sourceObject, targetBucket, targetObject string
	if req.ImageMediaId > 0 {
		asset, err := l.prepareImage(req.ImageMediaId, userID)
		if err != nil {
			return nil, err
		}
		sourceBucket, sourceObject = asset.Bucket, asset.ObjectKey
		targetBucket = l.svcCtx.Config.Minio.FormalBucket
		ext := strings.ToLower(filepath.Ext(asset.OriginName))
		if ext == "" {
			ext = ".jpg"
		}
		targetObject = fmt.Sprintf("comments/%d/%s%s", userID, uuid.NewString(), ext)
		if _, err = l.svcCtx.MinioClient.CopyObject(l.ctx, minio.CopyDestOptions{Bucket: targetBucket, Object: targetObject}, minio.CopySrcOptions{Bucket: sourceBucket, Object: sourceObject}); err != nil {
			return nil, fmt.Errorf("保存评论图片失败: %w", err)
		}
		input.ImageBucket, input.ImageObjectKey = targetBucket, targetObject
	}

	created, err := l.svcCtx.CommentRepo.Create(l.ctx, input)
	if err != nil {
		if targetObject != "" {
			_ = l.svcCtx.MinioClient.RemoveObject(context.Background(), targetBucket, targetObject, minio.RemoveObjectOptions{})
		}
		switch {
		case errors.Is(err, commentrepo.ErrWorkNotFound):
			return nil, fmt.Errorf("作品不存在或暂不可评论")
		case errors.Is(err, commentrepo.ErrParentNotFound):
			return nil, fmt.Errorf("回复的评论不存在")
		case errors.Is(err, commentrepo.ErrImageUnavailable):
			return nil, fmt.Errorf("评论图片状态已变化，请重新上传")
		default:
			l.Errorf("创建评论失败: %v", err)
			return nil, fmt.Errorf("评论发布失败")
		}
	}
	if sourceObject != "" {
		_ = l.svcCtx.MinioClient.RemoveObject(context.Background(), sourceBucket, sourceObject, minio.RemoveObjectOptions{})
	}
	row, err := l.svcCtx.CommentRepo.FindByID(l.ctx, int64(created.ID))
	if err != nil {
		l.Errorf("评论已创建但读取返回数据失败(commentId=%d): %v", created.ID, err)
		rootID := int64(created.ID)
		if created.RootID > 0 {
			rootID = created.RootID
		}
		fallback := types.CommentItem{
			CommentId: int64(created.ID), WorkId: created.WorkID, RootCommentId: rootID, ParentCommentId: created.ParentID,
			ReplyToUserId: created.ReplyToUserID, Content: created.Content, CreatedAt: created.CreatedAt.Format(time.RFC3339),
			ImageUrl: presignCommentImage(l.ctx, l.svcCtx, targetBucket, targetObject), Replies: make([]types.CommentItem, 0),
			Author: types.CommentAuthor{UserId: user.ID, NickName: normalizeName(user.Nickname), AvatarUrl: publicAvatarURL(l.svcCtx, user.Avatar)},
			Like:   types.RecommendLikeInfo{Liked: false, Count: 0, Text: "0"},
		}
		return &types.CreateCommentResp{Comment: fallback}, nil
	}
	item := buildCommentItem(l.ctx, l.svcCtx, *row)
	return &types.CreateCommentResp{Comment: item}, nil
}

func (l *CreateCommentLogic) prepareImage(mediaID, userID int64) (*mediarepo.MediaAsset, error) {
	asset, err := l.svcCtx.MediaRepo.FindAssetByID(l.ctx, mediaID)
	if err != nil {
		return nil, fmt.Errorf("评论图片不存在")
	}
	if asset.UserID != userID || asset.ResourceType != "image" || asset.Status != "uploaded" || asset.BoundCommentID != 0 || asset.BoundWorkID != 0 {
		return nil, fmt.Errorf("评论图片不可用")
	}
	if asset.FileSize <= 0 || asset.FileSize > 10<<20 {
		return nil, fmt.Errorf("评论图片不能超过10MB")
	}
	bucket := l.svcCtx.Config.Minio.FormalBucket
	if bucket == "" {
		return nil, fmt.Errorf("未配置正式资源桶")
	}
	exists, err := l.svcCtx.MinioClient.BucketExists(l.ctx, bucket)
	if err != nil {
		return nil, fmt.Errorf("检查评论图片存储失败: %w", err)
	}
	if !exists {
		if err = l.svcCtx.MinioClient.MakeBucket(l.ctx, bucket, minio.MakeBucketOptions{}); err != nil {
			return nil, fmt.Errorf("创建评论图片存储失败: %w", err)
		}
	}
	return asset, nil
}
