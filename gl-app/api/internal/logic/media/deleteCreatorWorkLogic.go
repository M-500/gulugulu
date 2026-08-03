package media

import (
	"context"
	"fmt"

	"gl-app/api/internal/svc"
	"gl-app/api/internal/types"
	"gl-app/pkg/ctxdata"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeleteCreatorWorkLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 删除当前用户作品
func NewDeleteCreatorWorkLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteCreatorWorkLogic {
	return &DeleteCreatorWorkLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DeleteCreatorWorkLogic) DeleteCreatorWork(req *types.WorkIdReq) (resp *types.CreatorWorkMutationResp, err error) {
	userID := ctxdata.GetUidFromCtx(l.ctx)
	// 作品采用软删除，正式媒体文件由后续资源清理任务统一回收。
	updated, err := l.svcCtx.WorkRepo.DeleteCreatorWork(l.ctx, req.WorkId, userID)
	if err != nil {
		return nil, fmt.Errorf("删除作品失败: %w", err)
	}
	if !updated {
		return nil, fmt.Errorf("作品不存在或已经删除")
	}

	return &types.CreatorWorkMutationResp{WorkId: req.WorkId, Updated: true}, nil
}
