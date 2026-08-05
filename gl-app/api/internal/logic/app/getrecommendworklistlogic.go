package app

import (
	"context"
	"fmt"
	"strings"
	"time"

	"gl-app/api/internal/constants"
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
	// 如果有传用户ID，那么就应该对特定用户推荐特定的作品
	hasMore := int64(len(rows)) > pageSize
	if hasMore {
		rows = rows[:pageSize]
	}
	resp := &types.RecommendWorkListResp{
		Page: page, PageSize: pageSize, HasMore: hasMore,
		List: make([]types.RecommendWorkItem, 0, len(rows)),
	}
	workIDs := make([]int64, 0, len(rows))
	for _, row := range rows {
		workIDs = append(workIDs, row.ID)
	}
	likes := queryWorkLikes(l.ctx, l.svcCtx, req.UserID, workIDs)
	for _, row := range rows {
		item := types.RecommendWorkItem{
			WorkId: row.ID, Type: row.Type, Title: row.Title,
			ContentExcerpt: strings.TrimSpace(row.ContentExcerpt), DurationMs: row.DurationMs,
			Author: types.RecommendAuthorItem{UserId: row.AuthorID, NickName: normalizeAuthorName(row.AuthorName),
				AvatarUrl: buildAppPublicAvatarURL(l.svcCtx, row.AuthorAvatar)},
			Like: likes[row.ID],
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

func queryWorkLikes(ctx context.Context, svcCtx *svc.ServiceContext, viewerUserID int64, workIDs []int64) map[int64]types.RecommendLikeInfo {
	result := make(map[int64]types.RecommendLikeInfo, len(workIDs))
	rows, err := svcCtx.InteractiveRepo.FindByResourceIDs(ctx, workIDs, constants.WorkType)
	likedStates := make(map[int64]bool)
	if viewerUserID > 0 {
		if values, likedErr := svcCtx.InteractiveRepo.FindUserLikedResourceIDs(ctx, viewerUserID, workIDs, constants.WorkType); likedErr == nil {
			likedStates = values
		}
	}
	for _, workID := range workIDs {
		count := int64(0)
		if err == nil && rows[workID] != nil {
			count = rows[workID].LikeCount
		}
		result[workID] = types.RecommendLikeInfo{Liked: likedStates[workID], Count: count, Text: formatLikeText(count)}
	}
	return result
}

func formatLikeText(count int64) string {
	if count >= 10000 {
		value := float64(count) / 10000
		return strings.TrimRight(strings.TrimRight(fmt.Sprintf("%.1f", value), "0"), ".") + "万"
	}
	return fmt.Sprintf("%d", count)
}
