package media

import (
	"context"
	"fmt"
	"strings"

	"gl-app/api/internal/svc"
	"gl-app/api/internal/types"
	"gl-app/pkg/ctxdata"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateCreatorWorkVisibilityLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 修改当前用户作品可见性
func NewUpdateCreatorWorkVisibilityLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateCreatorWorkVisibilityLogic {
	return &UpdateCreatorWorkVisibilityLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdateCreatorWorkVisibilityLogic) UpdateCreatorWorkVisibility(req *types.UpdateCreatorWorkVisibilityReq) (resp *types.CreatorWorkMutationResp, err error) {
	userID := ctxdata.GetUidFromCtx(l.ctx)
	visibility := strings.ToLower(strings.TrimSpace(req.Visibility))
	if visibility != "public" && visibility != "private" && visibility != "mutual" {
		return nil, fmt.Errorf("visibility必须是public、private或mutual")
	}
	updated, err := l.svcCtx.WorkRepo.UpdateCreatorVisibility(l.ctx, req.WorkId, userID, visibility)
	if err != nil {
		return nil, fmt.Errorf("修改作品权限失败: %w", err)
	}
	if !updated {
		return nil, fmt.Errorf("作品不存在或当前状态不可修改")
	}

	return &types.CreatorWorkMutationResp{WorkId: req.WorkId, Updated: true}, nil
}
