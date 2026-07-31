package media

import (
	"context"
	"fmt"

	"gl-app/api/internal/svc"
	"gl-app/api/internal/types"
	"gl-app/pkg/ctxdata"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetWorkPlaylistLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 动态生成作品HLS播放清单
func NewGetWorkPlaylistLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetWorkPlaylistLogic {
	return &GetWorkPlaylistLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetWorkPlaylistLogic) GetWorkPlaylist(req *types.WorkIdReq) (string, error) {
	userID := ctxdata.GetUidFromCtx(l.ctx)
	if userID <= 0 {
		return "", fmt.Errorf("用户未登录")
	}

	var asset struct {
		Bucket    string `db:"formal_bucket"`
		ObjectKey string `db:"formal_object_key"`
	}
	if err := l.svcCtx.SqlConn.QueryRowCtx(l.ctx, &asset, `SELECT ma.formal_bucket,ma.formal_object_key
		FROM work w
		JOIN work_asset wa ON wa.work_id=w.id AND wa.role='video'
		JOIN media_asset ma ON ma.id=wa.media_asset_id
		WHERE w.id=? AND w.user_id=? AND w.deleted_at IS NULL AND ma.status='ready'
		LIMIT 1`, req.WorkId, userID); err != nil {
		return "", fmt.Errorf("视频播放清单不存在: %w", err)
	}

	return BuildSignedHLSPlaylist(l.ctx, l.svcCtx, asset.Bucket, asset.ObjectKey)
}
