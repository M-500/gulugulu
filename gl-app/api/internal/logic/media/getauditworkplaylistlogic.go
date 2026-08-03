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

	asset, err := l.svcCtx.WorkRepo.FindVideoAsset(l.ctx, req.WorkId, 0, workrepo.VideoAccessAudit)
	if err != nil {
		return "", fmt.Errorf("审核视频播放清单不存在: %w", err)
	}

	return BuildSignedHLSPlaylist(l.ctx, l.svcCtx, asset.Bucket, asset.ObjectKey)
}
