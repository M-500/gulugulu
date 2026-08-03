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

type GetAuditWorkListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 查询审核中心作品列表
func NewGetAuditWorkListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetAuditWorkListLogic {
	return &GetAuditWorkListLogic{Logger: logx.WithContext(ctx), ctx: ctx, svcCtx: svcCtx}
}

func (l *GetAuditWorkListLogic) GetAuditWorkList(req *types.AuditWorkListReq) (*types.AuditWorkListResp, error) {
	if err := ensureAuditPermission(l.svcCtx, ctxdata.GetUidFromCtx(l.ctx)); err != nil {
		return nil, err
	}
	page, pageSize := normalizeCreatorWorkPage(req.Page, req.PageSize)
	status := strings.ToLower(strings.TrimSpace(req.Status))
	if status == "" {
		status = "pending"
	}
	if status != "all" && status != "pending" && status != "approved" && status != "rejected" {
		return nil, fmt.Errorf("status必须是all、pending、approved或rejected")
	}
	workType := strings.ToLower(strings.TrimSpace(req.Type))
	if workType != "" && workType != "image" && workType != "video" {
		return nil, fmt.Errorf("type必须是image或video")
	}
	summary, total, rows, err := l.svcCtx.WorkRepo.ListAuditWorks(l.ctx, workrepo.AuditListQuery{
		Status: status, WorkType: workType, Keyword: strings.TrimSpace(req.Keyword),
		Limit: pageSize, Offset: (page - 1) * pageSize,
	})
	if err != nil {
		return nil, fmt.Errorf("查询审核作品列表失败: %w", err)
	}
	resp := &types.AuditWorkListResp{
		Total: total, Page: page, PageSize: pageSize, PendingCount: summary.PendingCount,
		ApprovedCount: summary.ApprovedCount, RejectedCount: summary.RejectedCount,
		List: make([]types.AuditWorkListItem, 0, len(rows)),
	}
	for _, row := range rows {
		item := types.AuditWorkListItem{
			WorkId: row.ID, Type: row.Type, Title: row.Title, ContentExcerpt: row.ContentExcerpt,
			DurationMs: row.DurationMs, Visibility: row.Visibility, AuthorId: row.AuthorID,
			AuthorName: row.AuthorName, ReviewStatus: row.ReviewStatus, PublishStatus: row.PublishStatus,
			ReviewReason: row.ReviewReason, SubmittedAt: row.CreatedAt.Format(time.RFC3339),
			CoverUrl: l.presignAuditAsset(row.CoverBucket, row.CoverObjectKey),
		}
		if row.ReviewedAt != nil {
			item.ReviewedAt = row.ReviewedAt.Format(time.RFC3339)
		}
		resp.List = append(resp.List, item)
	}
	return resp, nil
}

func (l *GetAuditWorkListLogic) presignAuditAsset(bucket, objectKey string) string {
	if bucket == "" || objectKey == "" {
		return ""
	}
	value, err := l.svcCtx.MinioClient.PresignedGetObject(l.ctx, bucket, objectKey, 30*time.Minute, nil)
	if err != nil {
		l.Errorf("生成审核资源地址失败: %v", err)
		return ""
	}
	return value.String()
}
