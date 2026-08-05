package interactive

import (
	"net/http"

	interactivelogic "gl-app/api/internal/logic/interactive"
	"gl-app/api/internal/svc"
	"gl-app/api/internal/types"

	"github.com/zeromicro/go-zero/rest/httpx"
)

// UnlikeResourceHandler 取消点赞作品、评论或头像。
func UnlikeResourceHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.LikeResourceReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}
		resp, err := interactivelogic.NewMutateLikeLogic(r.Context(), svcCtx).Mutate(&req, false)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}
		httpx.OkJsonCtx(r.Context(), w, resp)
	}
}
