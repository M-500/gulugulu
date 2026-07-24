package media

import (
	"fmt"
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
	"gl-app/api/internal/logic/media"
	"gl-app/api/internal/svc"
	"gl-app/api/internal/types"
)

// 发布作品
func CreateWorkHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		r.Body = http.MaxBytesReader(w, r.Body, 6<<20)
		if err := r.ParseMultipartForm(12 << 20); err != nil {
			httpx.ErrorCtx(r.Context(), w, fmt.Errorf("解析发布表单失败: %w", err))
			return
		}
		var req types.CreateWorkReq
		req.Payload = r.FormValue("payload")
		cover, coverHeader, err := r.FormFile("cover")
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, fmt.Errorf("封面图不能为空: %w", err))
			return
		}
		defer cover.Close()

		l := media.NewCreateWorkLogic(r.Context(), svcCtx)
		resp, err := l.CreateWork(&req, cover, coverHeader, r.Header.Get("Idempotency-Key"))
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
