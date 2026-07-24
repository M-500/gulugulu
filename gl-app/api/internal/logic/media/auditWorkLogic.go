package media

import (
	"context"
	"database/sql"
	"fmt"
	"slices"
	"strings"
	"time"

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
	return &AuditWorkLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *AuditWorkLogic) AuditWork(req *types.AuditWorkReq) (resp *types.AuditWorkResp, err error) {
	reviewerID := ctxdata.GetUidFromCtx(l.ctx)
	if !slices.Contains(l.svcCtx.Config.Audit.AdminUserIds, reviewerID) {
		return nil, fmt.Errorf("当前用户没有审核权限")
	}
	decision := strings.ToLower(strings.TrimSpace(req.Decision))
	if decision != "approve" && decision != "reject" {
		return nil, fmt.Errorf("decision必须是approve或reject")
	}
	if decision == "reject" && strings.TrimSpace(req.Reason) == "" {
		return nil, fmt.Errorf("审核拒绝必须填写原因")
	}

	var row struct {
		ProcessStatus string       `db:"process_status"`
		ReviewStatus  string       `db:"review_status"`
		ScheduledAt   sql.NullTime `db:"scheduled_at"`
	}
	if err = l.svcCtx.SqlConn.QueryRowCtx(l.ctx, &row,
		"SELECT process_status,review_status,scheduled_at FROM work WHERE id=? AND deleted_at IS NULL",
		req.WorkId); err != nil {
		return nil, fmt.Errorf("作品不存在: %w", err)
	}
	if row.ProcessStatus != "succeeded" || row.ReviewStatus != "pending_review" {
		return nil, fmt.Errorf("作品当前状态不可审核")
	}

	reviewStatus := "rejected"
	publishStatus := "pending"
	var publishedAt any
	if decision == "approve" {
		reviewStatus = "approved"
		if row.ScheduledAt.Valid && row.ScheduledAt.Time.After(time.Now()) {
			publishStatus = "scheduled"
		} else {
			publishStatus = "published"
			publishedAt = time.Now().UTC()
		}
	}
	result, err := l.svcCtx.SqlConn.ExecCtx(l.ctx, `UPDATE work SET review_status=?,publish_status=?,reviewed_by=?,
		reviewed_at=NOW(3),review_reason=?,published_at=? WHERE id=? AND review_status='pending_review'`,
		reviewStatus, publishStatus, reviewerID, strings.TrimSpace(req.Reason), publishedAt, req.WorkId)
	if err != nil {
		return nil, err
	}
	affected, _ := result.RowsAffected()
	if affected != 1 {
		return nil, fmt.Errorf("作品已被其他审核员处理")
	}

	return &types.AuditWorkResp{WorkId: req.WorkId, ReviewStatus: reviewStatus, PublishStatus: publishStatus}, nil
}
