package media

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
	"gl-app/api/internal/logic/media"
	"gl-app/api/internal/svc"
	"gl-app/api/internal/types"
)

// 查询审核作品详情
func GetAuditWorkDetailHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.WorkIdReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := media.NewGetAuditWorkDetailLogic(r.Context(), svcCtx)
		resp, err := l.GetAuditWorkDetail(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
