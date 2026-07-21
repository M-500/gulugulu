package types

type CaptchaResponse struct {
	CaptchaID string `json:"captcha_id"`
	PicPath   string `json:"pic_path"`
}
