package media

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
	"gl-app/api/internal/logic/media"
	"gl-app/api/internal/svc"
	"gl-app/api/internal/types"
)

// 查询当前用户所有处理完成的作品
func GetCreatorWorkListHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.CreatorWorkListReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := media.NewGetCreatorWorkListLogic(r.Context(), svcCtx)
		resp, err := l.GetCreatorWorkList(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
