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

type UpdateCreatorWorkTitleLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 修改当前用户作品标题
func NewUpdateCreatorWorkTitleLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateCreatorWorkTitleLogic {
	return &UpdateCreatorWorkTitleLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdateCreatorWorkTitleLogic) UpdateCreatorWorkTitle(req *types.UpdateCreatorWorkTitleReq) (resp *types.CreatorWorkMutationResp, err error) {
	userID := ctxdata.GetUidFromCtx(l.ctx)
	title := strings.TrimSpace(req.Title)
	if title == "" || len([]rune(title)) > 50 {
		return nil, fmt.Errorf("作品标题不能为空且最多50个字符")
	}
	result, err := l.svcCtx.SqlConn.ExecCtx(l.ctx, `UPDATE work SET title=?
		WHERE id=? AND user_id=? AND process_status='succeeded' AND deleted_at IS NULL`,
		title, req.WorkId, userID)
	if err != nil {
		return nil, fmt.Errorf("修改作品标题失败: %w", err)
	}
	affected, _ := result.RowsAffected()
	if affected == 0 {
		return nil, fmt.Errorf("作品不存在或当前状态不可编辑")
	}

	return &types.CreatorWorkMutationResp{WorkId: req.WorkId, Updated: true}, nil
}
