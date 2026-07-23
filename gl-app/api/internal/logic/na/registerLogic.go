package na

import (
	"context"
	"database/sql"
	"errors"
	"gl-app/api/internal/models/user"
	"gl-app/pkg/utils/cryptx"
	"time"

	"gl-app/api/internal/svc"
	"gl-app/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type RegisterLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// register
func NewRegisterLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RegisterLogic {
	return &RegisterLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *RegisterLogic) Register(req *types.RegisterReq) (resp *types.RegisterResp, err error) {
	if req.Email == "" {
		return nil, errors.New("email不能为空")
	}
	_, err = l.svcCtx.UserRepo.FindOneByEmail(l.ctx, req.Email)
	if err == nil {
		return nil, errors.New("email 已经被占用，无法注册")
	}
	if !errors.Is(err, user.ErrNotFound) {
		return nil, errors.New("email 已经被占用，无法注册")
	}
	now := time.Now()
	// 校验邮箱验证码 TODO
	userInfo := user.User{
		CreatedAt: now,
		UpdatedAt: now,
		DeletedAt: sql.NullTime{
			Valid: false,
		},
		Email:    req.Email,
		Nickname: req.NickName,
		Password: cryptx.PasswordEncrypt(l.svcCtx.Config.Salt, req.Password),
		Sex:      0,
		LastLoginAt: sql.NullTime{
			Valid: false,
		},
	}
	insertResult, err := l.svcCtx.UserRepo.Insert(l.ctx, &userInfo)
	if err != nil {
		return nil, errors.New("注册失败")
	}
	userId, err := insertResult.LastInsertId()
	if err != nil {
		return nil, err
	}
	// 组装Token
	generateTokenLogic := NewGenerateTokenLogic(l.ctx, l.svcCtx)
	tokenResp, err := generateTokenLogic.GenerateToken(userId)
	if err != nil {
		return nil, errors.New("生成token失败")
	}
	res := types.RegisterResp{
		AccessToken:  tokenResp.AccessToken,
		AccessExpire: tokenResp.AccessExpire,
		RefreshAfter: tokenResp.RefreshAfter,
	}
	return &res, nil
}
