package user

import (
	"fmt"
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
	"gl-app/api/internal/logic/user"
	"gl-app/api/internal/svc"
)

// 上传当前用户头像
func UploadCurrentUserAvatarHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// 头像文件直接上传到后端，再由业务层写入MinIO公共桶。请求体大小由全局MaxBytes控制。
		if err := r.ParseMultipartForm(32 << 20); err != nil {
			httpx.ErrorCtx(r.Context(), w, fmt.Errorf("解析头像文件失败: %w", err))
			return
		}
		file, header, err := r.FormFile("avatar")
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, fmt.Errorf("请选择头像文件: %w", err))
			return
		}
		defer file.Close()

		l := user.NewUploadCurrentUserAvatarLogic(r.Context(), svcCtx)
		resp, err := l.UploadCurrentUserAvatar(file, header)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
