package comment

import (
	"context"
	"fmt"
	"net/url"
	"path"
	"strings"
	"time"

	commentrepo "gl-app/api/internal/repo/comment_repo"
	"gl-app/api/internal/svc"
	"gl-app/api/internal/types"
)

func normalizeCommentPage(page, pageSize int64) (int64, int64) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}
	if pageSize > 50 {
		pageSize = 50
	}
	return page, pageSize
}

func buildCommentItem(ctx context.Context, svcCtx *svc.ServiceContext, row commentrepo.CommentView) types.CommentItem {
	rootID := row.RootID
	if row.ParentID == 0 {
		rootID = row.ID
	}
	return types.CommentItem{
		CommentId: row.ID, WorkId: row.WorkID, RootCommentId: rootID,
		ParentCommentId: row.ParentID, ReplyToUserId: row.ReplyToUserID, ReplyToName: row.ReplyToName,
		Content: row.Content, ImageUrl: presignCommentImage(ctx, svcCtx, row.ImageBucket, row.ImageObjectKey),
		CreatedAt: row.CreatedAt.Format(time.RFC3339), ReplyCount: row.ReplyCount,
		ReplyHasMore: false, Replies: make([]types.CommentItem, 0),
		Author: types.CommentAuthor{UserId: row.UserID, NickName: normalizeName(row.AuthorName), AvatarUrl: publicAvatarURL(svcCtx, row.AuthorAvatar)},
		Like:   types.RecommendLikeInfo{Liked: false, Count: row.LikeCount, Text: formatCount(row.LikeCount)},
	}
}

func presignCommentImage(ctx context.Context, svcCtx *svc.ServiceContext, bucket, objectKey string) string {
	if bucket == "" || objectKey == "" {
		return ""
	}
	value, err := svcCtx.MinioClient.PresignedGetObject(ctx, bucket, objectKey, 30*time.Minute, nil)
	if err != nil {
		return ""
	}
	return value.String()
}

func publicAvatarURL(svcCtx *svc.ServiceContext, objectKey string) string {
	objectKey = strings.TrimSpace(objectKey)
	if objectKey == "" {
		return ""
	}
	bucket := svcCtx.Config.Minio.PublicBucket
	if bucket == "" {
		bucket = svcCtx.Config.Minio.FormalBucket
	}
	scheme := "http"
	if svcCtx.Config.Minio.UseSSL {
		scheme = "https"
	}
	endpoint := strings.TrimPrefix(strings.TrimPrefix(svcCtx.Config.Minio.Endpoint, "http://"), "https://")
	return (&url.URL{Scheme: scheme, Host: endpoint, Path: path.Join(bucket, objectKey)}).String()
}

func normalizeName(value string) string {
	value = strings.TrimSpace(value)
	if value == "" {
		return "咕噜用户"
	}
	return value
}

func formatCount(count int64) string {
	if count >= 10000 {
		return fmt.Sprintf("%.1f万", float64(count)/10000)
	}
	return fmt.Sprintf("%d", count)
}
