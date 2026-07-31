package user

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
	"gl-app/api/internal/logic/user"
	"gl-app/api/internal/svc"
)

// 修改当前用户基本信息
func UpdateUserInfoHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := r.ParseMultipartForm(32 << 20); err != nil {
			httpx.ErrorCtx(r.Context(), w, fmt.Errorf("解析用户资料表单失败: %w", err))
			return
		}
		nickname := r.FormValue("nickName")
		file, header, err := r.FormFile("avatar")
		if err != nil && !errors.Is(err, http.ErrMissingFile) {
			httpx.ErrorCtx(r.Context(), w, fmt.Errorf("读取头像文件失败: %w", err))
			return
		}
		if file != nil {
			defer file.Close()
		}

		l := user.NewUpdateUserInfoLogic(r.Context(), svcCtx)
		resp, err := l.UpdateUserInfo(nickname, file, header)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
