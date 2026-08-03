package media

import (
	"context"
	"fmt"
	"strings"
	"time"

	workrepo "gl-app/api/internal/repo/work_repo"
	"gl-app/api/internal/svc"
	"gl-app/api/internal/types"
	"gl-app/pkg/ctxdata"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetCreatorWorkListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 查询当前用户所有处理完成的作品
func NewGetCreatorWorkListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetCreatorWorkListLogic {
	return &GetCreatorWorkListLogic{Logger: logx.WithContext(ctx), ctx: ctx, svcCtx: svcCtx}
}

func (l *GetCreatorWorkListLogic) GetCreatorWorkList(req *types.CreatorWorkListReq) (*types.CreatorWorkListResp, error) {
	userID := ctxdata.GetUidFromCtx(l.ctx)
	if userID <= 0 {
		return nil, fmt.Errorf("用户未登录")
	}
	page, pageSize := normalizeCreatorWorkPage(req.Page, req.PageSize)
	status := strings.ToLower(strings.TrimSpace(req.Status))
	if status == "" {
		status = "all"
	}
	if status != "all" && status != "published" && status != "reviewing" && status != "rejected" {
		return nil, fmt.Errorf("status必须是all、published、reviewing或rejected")
	}
	summary, total, rows, err := l.svcCtx.WorkRepo.ListCreatorWorks(l.ctx, workrepo.CreatorListQuery{
		UserID: userID, Status: status, Keyword: strings.TrimSpace(req.Keyword),
		Limit: pageSize, Offset: (page - 1) * pageSize,
	})
	if err != nil {
		return nil, fmt.Errorf("查询作品列表失败: %w", err)
	}
	resp := &types.CreatorWorkListResp{
		Total: total, Page: page, PageSize: pageSize, AllCount: summary.AllCount,
		PublishedCount: summary.PublishedCount, ReviewingCount: summary.ReviewingCount,
		RejectedCount: summary.RejectedCount, List: make([]types.CreatorWorkListItem, 0, len(rows)),
	}
	for _, row := range rows {
		item := types.CreatorWorkListItem{
			WorkId: row.ID, Type: row.Type, Title: row.Title, DurationMs: row.DurationMs,
			Visibility: row.Visibility, ReviewStatus: row.ReviewStatus, PublishStatus: row.PublishStatus,
			ReviewReason: row.ReviewReason, CreatedAt: row.CreatedAt.Format(time.RFC3339),
		}
		if row.ScheduledAt != nil {
			item.ScheduledAt = row.ScheduledAt.Format(time.RFC3339)
		}
		if row.PublishedAt != nil {
			item.PublishedAt = row.PublishedAt.Format(time.RFC3339)
		}
		item.CoverUrl = l.presignCreatorCover(row.ID, row.CoverBucket, row.CoverObjectKey)
		resp.List = append(resp.List, item)
	}
	return resp, nil
}

func (l *GetCreatorWorkListLogic) presignCreatorCover(workID int64, bucket, objectKey string) string {
	if bucket == "" || objectKey == "" {
		return ""
	}
	value, err := l.svcCtx.MinioClient.PresignedGetObject(l.ctx, bucket, objectKey, 15*time.Minute, nil)
	if err != nil {
		l.Errorf("生成作品%d封面地址失败: %v", workID, err)
		return ""
	}
	return value.String()
}

func normalizeCreatorWorkPage(page, pageSize int64) (int64, int64) {
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
