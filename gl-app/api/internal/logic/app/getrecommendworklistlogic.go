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

type GetRecommendWorkListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// App首页推荐作品列表
func NewGetRecommendWorkListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetRecommendWorkListLogic {
	return &GetRecommendWorkListLogic{Logger: logx.WithContext(ctx), ctx: ctx, svcCtx: svcCtx}
}

func (l *GetRecommendWorkListLogic) GetRecommendWorkList(req *types.RecommendWorkListReq) (*types.RecommendWorkListResp, error) {
	page, pageSize := normalizeRecommendPage(req.Page, req.PageSize)
	rows, err := l.svcCtx.WorkRepo.ListRecommend(l.ctx, pageSize+1, (page-1)*pageSize)
	if err != nil {
		l.Errorf("查询App推荐作品列表失败: %v", err)
		return nil, fmt.Errorf("查询推荐作品失败")
	}
	hasMore := int64(len(rows)) > pageSize
	if hasMore {
		rows = rows[:pageSize]
	}
	resp := &types.RecommendWorkListResp{
		Page: page, PageSize: pageSize, HasMore: hasMore,
		List: make([]types.RecommendWorkItem, 0, len(rows)),
	}
	for _, row := range rows {
		item := types.RecommendWorkItem{
			WorkId: row.ID, Type: row.Type, Title: row.Title,
			ContentExcerpt: strings.TrimSpace(row.ContentExcerpt), DurationMs: row.DurationMs,
			Author: types.RecommendAuthorItem{UserId: row.AuthorID, NickName: normalizeAuthorName(row.AuthorName),
				AvatarUrl: buildAppPublicAvatarURL(l.svcCtx, row.AuthorAvatar)},
			Like: mockRecommendLike(row.ID),
		}
		if row.PublishedAt != nil {
			item.PublishedAt = row.PublishedAt.Format(time.RFC3339)
		}
		item.CoverUrl = l.presignCover(row.ID, row.CoverBucket, row.CoverObjectKey)
		resp.List = append(resp.List, item)
	}
	return resp, nil
}

func (l *GetRecommendWorkListLogic) presignCover(workID int64, bucket, objectKey string) string {
	if bucket == "" || objectKey == "" {
		return ""
	}
	value, err := l.svcCtx.MinioClient.PresignedGetObject(l.ctx, bucket, objectKey, 30*time.Minute, nil)
	if err != nil {
		l.Errorf("生成作品%d封面预签名地址失败: %v", workID, err)
		return ""
	}
	return value.String()
}

func normalizeRecommendPage(page, pageSize int64) (int64, int64) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 20
	}
	if pageSize > 50 {
		pageSize = 50
	}
	return page, pageSize
}

func normalizeAuthorName(name string) string {
	name = strings.TrimSpace(name)
	if name == "" {
		return "咕噜用户"
	}
	return name
}

func mockRecommendLike(workID int64) types.RecommendLikeInfo {
	// 点赞表还没有接入，先用作品ID生成稳定假数据，避免每次刷新跳变。
	count := int64(37) + (workID*97)%23800
	return types.RecommendLikeInfo{Liked: false, Count: count, Text: formatLikeText(count)}
}

func formatLikeText(count int64) string {
	if count >= 10000 {
		value := float64(count) / 10000
		return strings.TrimRight(strings.TrimRight(fmt.Sprintf("%.1f", value), "0"), ".") + "万"
	}
	return fmt.Sprintf("%d", count)
}
