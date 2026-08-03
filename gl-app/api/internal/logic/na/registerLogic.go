package na

import (
	"context"
	"errors"
	"strings"

	userrepo "gl-app/api/internal/repo/user_repo"
	"gl-app/api/internal/svc"
	"gl-app/api/internal/types"
	"gl-app/pkg/utils/cryptx"

	"github.com/zeromicro/go-zero/core/logx"
)

type RegisterLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// register
func NewRegisterLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RegisterLogic {
	return &RegisterLogic{Logger: logx.WithContext(ctx), ctx: ctx, svcCtx: svcCtx}
}

func (l *RegisterLogic) Register(req *types.RegisterReq) (*types.RegisterResp, error) {
	email := strings.TrimSpace(req.Email)
	if email == "" {
		return nil, errors.New("email不能为空")
	}
	if _, err := l.svcCtx.UserRepo.FindOneByEmail(l.ctx, email); err == nil {
		return nil, errors.New("email 已经被占用，无法注册")
	} else if !errors.Is(err, userrepo.ErrNotFound) {
		return nil, errors.New("查询邮箱是否可用失败")
	}

	// 用户写入和自增ID回填均由 GORM Repo 完成。
	user := &userrepo.User{
		Email: email, Nickname: strings.TrimSpace(req.NickName),
		Password: cryptx.PasswordEncrypt(l.svcCtx.Config.Salt, req.Password), Sex: 0,
	}
	if err := l.svcCtx.UserRepo.Create(l.ctx, user); err != nil {
		return nil, errors.New("注册失败")
	}
	tokenResp, err := NewGenerateTokenLogic(l.ctx, l.svcCtx).GenerateToken(user.ID)
	if err != nil {
		return nil, errors.New("生成token失败")
	}
	return &types.RegisterResp{
		AccessToken: tokenResp.AccessToken, AccessExpire: tokenResp.AccessExpire, RefreshAfter: tokenResp.RefreshAfter,
	}, nil
}
