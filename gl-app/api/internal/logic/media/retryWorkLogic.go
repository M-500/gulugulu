package media

import (
	"context"
	"fmt"

	"gl-app/api/internal/svc"
	"gl-app/api/internal/types"
	"gl-app/pkg/ctxdata"

	"github.com/zeromicro/go-zero/core/logx"
)

type RetryWorkLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 重试失败的媒体处理任务
func NewRetryWorkLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RetryWorkLogic {
	return &RetryWorkLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *RetryWorkLogic) RetryWork(req *types.WorkIdReq) (resp *types.WorkStatusResp, err error) {
	userID := ctxdata.GetUidFromCtx(l.ctx)
	var workType string
	if err = l.svcCtx.SqlConn.QueryRowCtx(l.ctx, &workType,
		"SELECT type FROM work WHERE id=? AND user_id=? AND process_status='failed' AND deleted_at IS NULL",
		req.WorkId, userID); err != nil {
		return nil, fmt.Errorf("作品不存在或当前状态不可重试")
	}
	result, err := l.svcCtx.SqlConn.ExecCtx(l.ctx, `UPDATE work w JOIN media_process_task t ON t.work_id=w.id
		SET w.process_status='pending',w.review_status='waiting_process',t.status='pending',t.stage='queued',
		t.progress=0,t.error_message='',t.retry_count=t.retry_count+1
		WHERE w.id=? AND w.user_id=? AND w.process_status='failed'`, req.WorkId, userID)
	if err != nil {
		return nil, err
	}
	affected, _ := result.RowsAffected()
	if affected == 0 {
		return nil, fmt.Errorf("作品不存在或当前状态不可重试")
	}
	if workType == "image" {
		if err = promoteImageWork(l.ctx, l.svcCtx, req.WorkId); err != nil {
			markImagePromotionFailed(l.ctx, l.svcCtx, req.WorkId, err)
			return nil, err
		}
		return &types.WorkStatusResp{
			WorkId: req.WorkId, ProcessStatus: "succeeded", ReviewStatus: "pending_review",
			PublishStatus: "pending", Stage: "waiting_review", Progress: 100, Message: "图片迁移完成，等待审核",
		}, nil
	}
	if err = l.svcCtx.MediaQueue.PublishProcessWork(l.ctx, req.WorkId, 1); err != nil {
		return nil, err
	}

	return &types.WorkStatusResp{
		WorkId: req.WorkId, ProcessStatus: "pending", ReviewStatus: "waiting_process",
		PublishStatus: "pending", Stage: "queued", Progress: 0, Message: "已重新加入处理队列",
	}, nil
}
