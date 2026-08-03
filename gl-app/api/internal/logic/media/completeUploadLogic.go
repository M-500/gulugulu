package media

import (
	"context"
	"fmt"
	"time"

	"github.com/minio/minio-go/v7"
	"gl-app/api/internal/svc"
	"gl-app/api/internal/types"
	"gl-app/pkg/ctxdata"

	"github.com/zeromicro/go-zero/core/logx"
)

type CompleteUploadLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// complete direct upload
func NewCompleteUploadLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CompleteUploadLogic {
	return &CompleteUploadLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CompleteUploadLogic) CompleteUpload(req *types.CompleteUploadReq) (resp *types.CompleteUploadResp, err error) {
	userId := ctxdata.GetUidFromCtx(l.ctx)
	if userId <= 0 {
		return nil, fmt.Errorf("用户未登录")
	}

	asset, err := l.svcCtx.MediaRepo.FindAssetByID(l.ctx, req.MediaId)
	if err != nil {
		return nil, fmt.Errorf("媒体素材不存在: %w", err)
	}

	if asset.UserID != userId {
		return nil, fmt.Errorf("无权操作该媒体素材")
	}

	if asset.ObjectKey != req.ObjectKey {
		return nil, fmt.Errorf("objectKey不匹配")
	}

	info, err := l.svcCtx.MinioClient.StatObject(l.ctx, asset.Bucket, asset.ObjectKey, minio.StatObjectOptions{})
	if err != nil {
		return nil, fmt.Errorf("对象存储中未找到上传文件: %w", err)
	}

	asset.FileSize = info.Size
	asset.Status = "uploaded"
	if info.ContentType != "" {
		asset.ContentType = info.ContentType
	}

	if err := l.svcCtx.MediaRepo.UpdateAsset(l.ctx, asset); err != nil {
		return nil, fmt.Errorf("更新媒体素材状态失败: %w", err)
	}
	expiresIn := l.svcCtx.Config.Minio.PresignExpire
	if expiresIn <= 0 {
		expiresIn = 900
	}
	previewUrl, err := l.svcCtx.MinioClient.PresignedGetObject(l.ctx, asset.Bucket, asset.ObjectKey, time.Duration(expiresIn)*time.Second, nil)
	if err != nil {
		return nil, fmt.Errorf("生成预览地址失败: %w", err)
	}

	return &types.CompleteUploadResp{
		MediaId:    asset.ID,
		Status:     asset.Status,
		Bucket:     asset.Bucket,
		ObjectKey:  asset.ObjectKey,
		PreviewUrl: previewUrl.String(),
	}, nil
}
