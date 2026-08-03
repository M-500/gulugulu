package na

import (
	"context"
	"github.com/pkg/errors"
	userlogic "gl-app/api/internal/logic/user"
	userrepo "gl-app/api/internal/repo/user_repo"
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
		if errors.Is(err, userrepo.ErrNotFound) {
			return nil, errors.New("邮箱不存在")
		}
		return nil, errors.New("查询用户信息失败")
	}
	generateTokenLogic := NewGenerateTokenLogic(l.ctx, l.svcCtx)
	tokenResp, err := generateTokenLogic.GenerateToken(userInfo.ID)
	if err != nil {
		return nil, errors.New("生成token失败")
	}
	profile, err := userlogic.QueryUserProfile(l.ctx, l.svcCtx, userInfo.ID)
	if err != nil {
		l.Errorf("登录成功后查询用户基础信息失败，用户ID：%d，错误：%v", userInfo.ID, err)
		return nil, errors.New("查询用户基础信息失败")
	}
	res := types.LoginResp{
		AccessToken:  tokenResp.AccessToken,
		AccessExpire: tokenResp.AccessExpire,
		RefreshAfter: tokenResp.RefreshAfter,
		UserId:       profile.UserId,
		NickName:     profile.NickName,
		AvatarUrl:    profile.AvatarUrl,
	}
	return &res, nil
}
