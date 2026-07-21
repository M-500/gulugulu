package na

import (
	"context"

	"gl-app/api/internal/svc"
	"gl-app/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type CaptchaLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext

	imgSvc svc.ImgCaptchaSvc
}

// captcha
func NewCaptchaLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CaptchaLogic {
	return &CaptchaLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
		imgSvc: svc.NewImgCaptchaSvc(ctx),
	}
}

func (l *CaptchaLogic) Captcha() (resp *types.CaptchaResp, err error) {
	captcha, err := l.imgSvc.GetCaptcha()
	if err != nil {
		return nil, err
	}
	resp.CaptchaID = captcha.CaptchaID
	resp.CaptchaPath = captcha.PicPath
	return resp, err
}
