package media

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
	"gl-app/api/internal/logic/media"
	"gl-app/api/internal/svc"
	"gl-app/api/internal/types"
)

// 查询审核中心作品列表
func GetAuditWorkListHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.AuditWorkListReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := media.NewGetAuditWorkListLogic(r.Context(), svcCtx)
		resp, err := l.GetAuditWorkList(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
