package app

import (
	"context"
	"fmt"

	mediaLogic "gl-app/api/internal/logic/media"
	"gl-app/api/internal/svc"
	"gl-app/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetAppWorkPlaylistLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 获取App公开视频HLS播放清单
func NewGetAppWorkPlaylistLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetAppWorkPlaylistLogic {
	return &GetAppWorkPlaylistLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetAppWorkPlaylistLogic) GetAppWorkPlaylist(req *types.AppWorkDetailReq) (string, error) {
	if req == nil || req.WorkId <= 0 {
		return "", fmt.Errorf("作品ID不正确")
	}

	var asset struct {
		Bucket    string `db:"formal_bucket"`
		ObjectKey string `db:"formal_object_key"`
	}
	query := `SELECT ma.formal_bucket,ma.formal_object_key
		FROM work w
		INNER JOIN work_asset wa ON wa.work_id=w.id AND wa.role='video'
		INNER JOIN media_asset ma ON ma.id=wa.media_asset_id AND ma.deleted_at IS NULL
		WHERE w.id=?
			AND w.type='video'
			AND w.deleted_at IS NULL
			AND w.process_status='succeeded'
			AND w.review_status='approved'
			AND w.publish_status='published'
			AND w.visibility='public'
			AND ma.status='ready'
		LIMIT 1`
	if err := l.svcCtx.SqlConn.QueryRowCtx(l.ctx, &asset, query, req.WorkId); err != nil {
		l.Errorf("查询App公开视频%d播放清单失败: %v", req.WorkId, err)
		return "", fmt.Errorf("视频播放清单不存在")
	}

	return mediaLogic.BuildSignedHLSPlaylist(l.ctx, l.svcCtx, asset.Bucket, asset.ObjectKey)
}
