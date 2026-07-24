package media

import (
	"fmt"
	"slices"

	"gl-app/api/internal/svc"
)

func ensureAuditPermission(svcCtx *svc.ServiceContext, userID int64) error {
	if userID <= 0 || !slices.Contains(svcCtx.Config.Audit.AdminUserIds, userID) {
		return fmt.Errorf("当前用户没有审核权限")
	}
	return nil
}
