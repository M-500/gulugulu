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

type AuditWorkLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 审核作品
func NewAuditWorkLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AuditWorkLogic {
	return &AuditWorkLogic{Logger: logx.WithContext(ctx), ctx: ctx, svcCtx: svcCtx}
}

func (l *AuditWorkLogic) AuditWork(req *types.AuditWorkReq) (*types.AuditWorkResp, error) {
	reviewerID := ctxdata.GetUidFromCtx(l.ctx)
	if err := ensureAuditPermission(l.svcCtx, reviewerID); err != nil {
		return nil, err
	}
	decision := strings.ToLower(strings.TrimSpace(req.Decision))
	if decision != "approve" && decision != "reject" {
		return nil, fmt.Errorf("decision必须是approve或reject")
	}
	if decision == "reject" && strings.TrimSpace(req.Reason) == "" {
		return nil, fmt.Errorf("审核拒绝必须填写原因")
	}
	state, err := l.svcCtx.WorkRepo.FindAuditState(l.ctx, req.WorkId)
	if err != nil {
		return nil, fmt.Errorf("作品不存在: %w", err)
	}
	if state.ProcessStatus != "succeeded" || state.ReviewStatus != "pending_review" {
		return nil, fmt.Errorf("作品当前状态不可审核")
	}

	reviewStatus, publishStatus := "rejected", "pending"
	var publishedAt *time.Time
	if decision == "approve" {
		reviewStatus = "approved"
		if state.ScheduledAt != nil && state.ScheduledAt.After(time.Now()) {
			publishStatus = "scheduled"
		} else {
			publishStatus = "published"
			now := time.Now().UTC()
			publishedAt = &now
		}
	}
	updated, err := l.svcCtx.WorkRepo.Review(l.ctx, workrepo.ReviewInput{
		WorkID: req.WorkId, ReviewerID: reviewerID, ReviewStatus: reviewStatus,
		PublishStatus: publishStatus, Reason: strings.TrimSpace(req.Reason), PublishedAt: publishedAt,
	})
	if err != nil {
		return nil, err
	}
	if !updated {
		return nil, fmt.Errorf("作品已被其他审核员处理")
	}
	return &types.AuditWorkResp{WorkId: req.WorkId, ReviewStatus: reviewStatus, PublishStatus: publishStatus}, nil
}
