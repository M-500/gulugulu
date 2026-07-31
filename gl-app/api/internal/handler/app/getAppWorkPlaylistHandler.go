package app

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
	"gl-app/api/internal/logic/app"
	"gl-app/api/internal/svc"
	"gl-app/api/internal/types"
)

// 获取App公开视频HLS播放清单
func GetAppWorkPlaylistHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.AppWorkDetailReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := app.NewGetAppWorkPlaylistLogic(r.Context(), svcCtx)
		playlist, err := l.GetAppWorkPlaylist(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			w.Header().Set("Content-Type", "application/vnd.apple.mpegurl; charset=utf-8")
			w.Header().Set("Cache-Control", "public, max-age=60")
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte(playlist))
		}
	}
}
