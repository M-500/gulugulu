package media

import (
	"context"
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
	return &GetWorkStatusLogic{Logger: logx.WithContext(ctx), ctx: ctx, svcCtx: svcCtx}
}

func (l *GetWorkStatusLogic) GetWorkStatus(req *types.WorkIdReq) (*types.WorkStatusResp, error) {
	row, err := l.svcCtx.WorkRepo.FindWorkStatus(l.ctx, req.WorkId, ctxdata.GetUidFromCtx(l.ctx))
	if err != nil {
		return nil, fmt.Errorf("作品不存在: %w", err)
	}
	resp := &types.WorkStatusResp{
		WorkId: row.ID, ProcessStatus: row.ProcessStatus, ReviewStatus: row.ReviewStatus,
		PublishStatus: row.PublishStatus, Stage: row.Stage, Progress: row.Progress,
		Message: statusMessage(row.ProcessStatus, row.ReviewStatus, row.PublishStatus), FailureReason: row.ErrorMessage,
	}
	if row.ScheduledAt != nil {
		resp.ScheduledAt = row.ScheduledAt.Format(time.RFC3339)
	}
	if row.PublishedAt != nil {
		resp.PublishedAt = row.PublishedAt.Format(time.RFC3339)
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
