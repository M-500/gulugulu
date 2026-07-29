package media

import (
	"context"
	"fmt"

	"gl-app/api/internal/svc"
	"gl-app/api/internal/types"
	"gl-app/pkg/ctxdata"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetAuditWorkPlaylistLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 动态生成审核作品HLS播放清单
func NewGetAuditWorkPlaylistLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetAuditWorkPlaylistLogic {
	return &GetAuditWorkPlaylistLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetAuditWorkPlaylistLogic) GetAuditWorkPlaylist(req *types.WorkIdReq) (string, error) {
	if err := ensureAuditPermission(l.svcCtx, ctxdata.GetUidFromCtx(l.ctx)); err != nil {
		return "", err
	}

	var asset struct {
		Bucket    string `db:"formal_bucket"`
		ObjectKey string `db:"formal_object_key"`
	}
	if err := l.svcCtx.SqlConn.QueryRowCtx(l.ctx, &asset, `SELECT ma.formal_bucket,ma.formal_object_key
		FROM work w
		JOIN work_asset wa ON wa.work_id=w.id AND wa.role='video'
		JOIN media_asset ma ON ma.id=wa.media_asset_id
		WHERE w.id=? AND w.process_status='succeeded' AND w.deleted_at IS NULL AND ma.status='ready'
		LIMIT 1`, req.WorkId); err != nil {
		return "", fmt.Errorf("审核视频播放清单不存在: %w", err)
	}

	return buildSignedHLSPlaylist(l.ctx, l.svcCtx, asset.Bucket, asset.ObjectKey)
}
