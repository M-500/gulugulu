package app

import (
	"context"
	"fmt"

	mediaLogic "gl-app/api/internal/logic/media"
	workrepo "gl-app/api/internal/repo/work_repo"
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

	asset, err := l.svcCtx.WorkRepo.FindVideoAsset(l.ctx, req.WorkId, 0, workrepo.VideoAccessPublic)
	if err != nil {
		l.Errorf("查询App公开视频%d播放清单失败: %v", req.WorkId, err)
		return "", fmt.Errorf("视频播放清单不存在")
	}

	return mediaLogic.BuildSignedHLSPlaylist(l.ctx, l.svcCtx, asset.Bucket, asset.ObjectKey)
}
