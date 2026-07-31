package user

import (
	"context"
	"fmt"
	"strings"
	"unicode/utf8"

	"gl-app/api/internal/svc"
	"gl-app/api/internal/types"
	"gl-app/pkg/ctxdata"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateCurrentUserLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 修改当前用户资料
func NewUpdateCurrentUserLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateCurrentUserLogic {
	return &UpdateCurrentUserLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdateCurrentUserLogic) UpdateCurrentUser(req *types.UpdateUserProfileReq) (resp *types.UserProfileResp, err error) {
	userID := ctxdata.GetUidFromCtx(l.ctx)
	nickname := strings.TrimSpace(req.NickName)
	length := utf8.RuneCountInString(nickname)
	if length < 2 || length > 32 {
		return nil, fmt.Errorf("昵称长度必须是2到32个字符")
	}
	if _, err = l.svcCtx.SqlConn.ExecCtx(l.ctx,
		"UPDATE user SET nickname=? WHERE id=? AND deleted_at IS NULL", nickname, userID); err != nil {
		return nil, fmt.Errorf("更新用户昵称失败: %w", err)
	}
	return QueryUserProfile(l.ctx, l.svcCtx, userID)
}
