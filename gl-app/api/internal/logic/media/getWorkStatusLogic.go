package media

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"gl-app/api/internal/svc"
	"gl-app/api/internal/types"
	"gl-app/pkg/ctxdata"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetWorkStatusLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 查询作品媒体处理和发布状态
func NewGetWorkStatusLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetWorkStatusLogic {
	return &GetWorkStatusLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetWorkStatusLogic) GetWorkStatus(req *types.WorkIdReq) (resp *types.WorkStatusResp, err error) {
	userID := ctxdata.GetUidFromCtx(l.ctx)
	var row struct {
		ID            int64        `db:"id"`
		ProcessStatus string       `db:"process_status"`
		ReviewStatus  string       `db:"review_status"`
		PublishStatus string       `db:"publish_status"`
		ScheduledAt   sql.NullTime `db:"scheduled_at"`
		PublishedAt   sql.NullTime `db:"published_at"`
		Stage         string       `db:"stage"`
		Progress      int64        `db:"progress"`
		ErrorMessage  string       `db:"error_message"`
	}
	err = l.svcCtx.SqlConn.QueryRowCtx(l.ctx, &row, `SELECT w.id,w.process_status,w.review_status,w.publish_status,
		w.scheduled_at,w.published_at,COALESCE(t.stage,'') AS stage,COALESCE(t.progress,0) AS progress,
		COALESCE(t.error_message,'') AS error_message
		FROM work w LEFT JOIN media_process_task t ON t.work_id=w.id AND t.task_type='process_work'
		WHERE w.id=? AND w.user_id=? AND w.deleted_at IS NULL LIMIT 1`, req.WorkId, userID)
	if err != nil {
		return nil, fmt.Errorf("作品不存在: %w", err)
	}

	resp = &types.WorkStatusResp{
		WorkId:        row.ID,
		ProcessStatus: row.ProcessStatus,
		ReviewStatus:  row.ReviewStatus,
		PublishStatus: row.PublishStatus,
		Stage:         row.Stage,
		Progress:      row.Progress,
		Message:       statusMessage(row.ProcessStatus, row.ReviewStatus, row.PublishStatus),
		FailureReason: row.ErrorMessage,
	}
	if row.ScheduledAt.Valid {
		resp.ScheduledAt = row.ScheduledAt.Time.Format(time.RFC3339)
	}
	if row.PublishedAt.Valid {
		resp.PublishedAt = row.PublishedAt.Time.Format(time.RFC3339)
	}
	return resp, nil
}

func statusMessage(processStatus, reviewStatus, publishStatus string) string {
	switch {
	case processStatus == "failed":
		return "素材处理失败"
	case processStatus != "succeeded":
		return "素材处理中"
	case reviewStatus == "pending_review":
		return "等待审核"
	case reviewStatus == "rejected":
		return "审核未通过"
	case publishStatus == "scheduled":
		return "审核通过，等待定时发布"
	case publishStatus == "published":
		return "已发布"
	default:
		return "处理中"
	}
}
