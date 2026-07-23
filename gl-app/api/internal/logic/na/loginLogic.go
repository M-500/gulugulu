package na

import (
	"context"
	"github.com/pkg/errors"
	"gl-app/api/internal/models/user"
	"strings"

	"gl-app/api/internal/svc"
	"gl-app/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type LoginLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
	imgSvc svc.ImgCaptchaSvc
}

// login
func NewLoginLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LoginLogic {
	return &LoginLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
		imgSvc: svc.NewImgCaptchaSvc(ctx),
	}
}

func (l *LoginLogic) Login(req *types.LoginReq) (resp *types.LoginResp, err error) {
	captcha := l.imgSvc.VerifyCaptcha(req.CaptchaID, req.CaptchaCode, true)
	if !captcha {
		return nil, errors.New("验证码错误")
	}
	req.Email = strings.TrimSpace(req.Email)
	if req.Email == "" {
		return nil, errors.New("邮箱不能为空")
	}
	userInfo, err := l.svcCtx.UserRepo.FindOneByEmail(l.ctx, req.Email)
	if err != nil {
		if errors.Is(err, user.ErrNotFound) {
			return nil, errors.New("邮箱不存在")
		}
		return nil, errors.New("查询用户信息失败")
	}
	generateTokenLogic := NewGenerateTokenLogic(l.ctx, l.svcCtx)
	tokenResp, err := generateTokenLogic.GenerateToken(userInfo.Id)
	if err != nil {
		return nil, errors.New("生成token失败")
	}
	res := types.LoginResp{
		AccessToken:  tokenResp.AccessToken,
		AccessExpire: tokenResp.AccessExpire,
		RefreshAfter: tokenResp.RefreshAfter,
	}
	return &res, nil
}
