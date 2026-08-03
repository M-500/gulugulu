package media

import (
	"context"
	"fmt"

	workrepo "gl-app/api/internal/repo/work_repo"
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

	asset, err := l.svcCtx.WorkRepo.FindVideoAsset(l.ctx, req.WorkId, userID, workrepo.VideoAccessOwner)
	if err != nil {
		return "", fmt.Errorf("视频播放清单不存在: %w", err)
	}

	return BuildSignedHLSPlaylist(l.ctx, l.svcCtx, asset.Bucket, asset.ObjectKey)
}
