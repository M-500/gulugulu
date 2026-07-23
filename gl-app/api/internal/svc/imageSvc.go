package svc

import (
	"context"
	"fmt"
	"github.com/mojocn/base64Captcha"
	"github.com/zeromicro/go-zero/core/logx"
	"gl-app/api/internal/types"
)

type ImgCaptchaSvc interface {
	// VerifyCaptcha verifies the captcha
	VerifyCaptcha(id, code string, clear bool) bool
	// GetCaptcha gets the captcha
	GetCaptcha() (*types.CaptchaResponse, error)
}

type imgCaptchaSvc struct {
	store base64Captcha.Store
	logx.Logger
}

func NewImgCaptchaSvc(ctx context.Context) ImgCaptchaSvc {
	return &imgCaptchaSvc{
		store:  base64Captcha.DefaultMemStore,
		Logger: logx.WithContext(ctx),
	}
}

func (i *imgCaptchaSvc) VerifyCaptcha(id, code string, clear bool) bool {
	return i.store.Verify(id, code, clear)
}

func (i *imgCaptchaSvc) GetCaptcha() (*types.CaptchaResponse, error) {
	driver := base64Captcha.NewDriverDigit(
		120, // height of png in pixels
		240, // width of png in pixels
		5,   // default number of captcha
		0.7, // 单个数字的最大偏斜因子
		130)
	cp := base64Captcha.NewCaptcha(driver, i.store)
	id, b64s, _, err := cp.Generate()
	if err != nil {
		i.Logger.Error("生成图片验证码失败", err)
		return nil, fmt.Errorf("生成图片验证码失败: %w", err)
	}
	res := types.CaptchaResponse{
		CaptchaID: id,
		PicPath:   b64s,
	}
	return &res, nil
}
