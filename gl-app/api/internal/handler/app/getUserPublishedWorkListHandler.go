package app

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
	"gl-app/api/internal/logic/app"
	"gl-app/api/internal/svc"
	"gl-app/api/internal/types"
)

// 分页获取用户已发布的作品；Bearer Token 可选，本人可查看非公开作品。
func GetUserPublishedWorkListHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.UserPublishedWorkListReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		viewerUserID := optionalJWTUserID(r, svcCtx.Config.JwtAuth.AccessSecret)
		l := app.NewGetUserPublishedWorkListLogic(r.Context(), svcCtx)
		resp, err := l.GetUserPublishedWorkList(&req, viewerUserID)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
