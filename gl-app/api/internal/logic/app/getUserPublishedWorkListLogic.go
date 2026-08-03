package app

import (
	"context"
	"fmt"
	"strings"
	"time"

	"gl-app/api/internal/svc"
	"gl-app/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetUserPublishedWorkListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 分页获取用户已发布的公开作品
func NewGetUserPublishedWorkListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetUserPublishedWorkListLogic {
	return &GetUserPublishedWorkListLogic{Logger: logx.WithContext(ctx), ctx: ctx, svcCtx: svcCtx}
}

func (l *GetUserPublishedWorkListLogic) GetUserPublishedWorkList(req *types.UserPublishedWorkListReq) (*types.UserPublishedWorkListResp, error) {
	if req == nil || req.UserId <= 0 {
		return nil, fmt.Errorf("用户ID不正确")
	}
	user, err := findAppPublicUser(l.ctx, l.svcCtx, req.UserId)
	if err != nil {
		return nil, err
	}
	page, pageSize := normalizeRecommendPage(req.Page, req.PageSize)
	offset := (page - 1) * pageSize
	total, rows, err := l.svcCtx.WorkRepo.ListPublishedByUser(l.ctx, req.UserId, pageSize, offset)
	if err != nil {
		l.Errorf("查询用户%d已发布作品失败: %v", req.UserId, err)
		return nil, fmt.Errorf("查询用户作品失败")
	}
	resp := &types.UserPublishedWorkListResp{
		Total: total, Page: page, PageSize: pageSize, HasMore: offset+int64(len(rows)) < total,
		List: make([]types.RecommendWorkItem, 0, len(rows)),
	}
	author := types.RecommendAuthorItem{
		UserId: user.ID, NickName: normalizeAuthorName(user.Nickname),
		AvatarUrl: buildAppPublicAvatarURL(l.svcCtx, user.Avatar),
	}
	for _, row := range rows {
		item := types.RecommendWorkItem{
			WorkId: row.ID, Type: row.Type, Title: row.Title,
			ContentExcerpt: strings.TrimSpace(row.ContentExcerpt), DurationMs: row.DurationMs,
			Author: author, Like: mockRecommendLike(row.ID),
		}
		if row.PublishedAt != nil {
			item.PublishedAt = row.PublishedAt.Format(time.RFC3339)
		}
		item.CoverUrl = l.presignUserWorkCover(row.ID, row.CoverBucket, row.CoverObjectKey)
		resp.List = append(resp.List, item)
	}
	return resp, nil
}

func (l *GetUserPublishedWorkListLogic) presignUserWorkCover(workID int64, bucket, objectKey string) string {
	if bucket == "" || objectKey == "" {
		return ""
	}
	value, err := l.svcCtx.MinioClient.PresignedGetObject(l.ctx, bucket, objectKey, 30*time.Minute, nil)
	if err != nil {
		l.Errorf("生成用户作品%d封面地址失败: %v", workID, err)
		return ""
	}
	return value.String()
}
