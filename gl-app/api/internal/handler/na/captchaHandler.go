package na

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
	"gl-app/api/internal/logic/na"
	"gl-app/api/internal/svc"
)

// captcha
func CaptchaHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := na.NewCaptchaLogic(r.Context(), svcCtx)
		resp, err := l.Captcha()
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
