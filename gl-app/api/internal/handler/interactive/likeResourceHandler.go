package interactive

import (
	"net/http"

	interactivelogic "gl-app/api/internal/logic/interactive"
	"gl-app/api/internal/svc"
	"gl-app/api/internal/types"

	"github.com/zeromicro/go-zero/rest/httpx"
)

// LikeResourceHandler 点赞作品、评论或头像。
func LikeResourceHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.LikeResourceReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}
		resp, err := interactivelogic.NewMutateLikeLogic(r.Context(), svcCtx).Mutate(&req, true)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}
		httpx.OkJsonCtx(r.Context(), w, resp)
	}
}
