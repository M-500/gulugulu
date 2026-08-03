package media

import (
	"context"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/minio/minio-go/v7"
	mediarepo "gl-app/api/internal/repo/media_repo"
	"gl-app/api/internal/svc"
)

type promotedImage struct {
	mediarepo.ProcessAsset
	targetObjectKey string
}

// promoteImageWork 在MinIO服务端直接复制图片，不需要下载文件，也不经过Kafka。
// 全部图片和封面复制成功并写入数据库后，作品立即进入待审核状态。
func promoteImageWork(ctx context.Context, svcCtx *svc.ServiceContext, workID int64) error {
	if err := ensureFormalMediaBucket(ctx, svcCtx); err != nil {
		return err
	}
	if err := svcCtx.MediaRepo.StartImagePromotion(ctx, workID); err != nil {
		return fmt.Errorf("更新图片处理状态失败: %w", err)
	}

	assets, err := svcCtx.MediaRepo.ListImageAssets(ctx, workID)
	if err != nil {
		return fmt.Errorf("查询图片作品素材失败: %w", err)
	}
	if len(assets) < 2 {
		return fmt.Errorf("图片作品素材不完整")
	}

	promoted := make([]promotedImage, 0, len(assets))
	for _, asset := range assets {
		if asset.Status == "ready" && asset.FormalObjectKey != "" {
			promoted = append(promoted, promotedImage{
				ProcessAsset:    asset,
				targetObjectKey: asset.FormalObjectKey,
			})
			continue
		}
		extension := strings.ToLower(filepath.Ext(asset.OriginName))
		if extension == "" {
			extension = ".jpg"
		}
		targetObjectKey := fmt.Sprintf("works/%d/%s/%d%s", workID, asset.Role, asset.ID, extension)
		source := minio.CopySrcOptions{Bucket: asset.Bucket, Object: asset.ObjectKey}
		destination := minio.CopyDestOptions{
			Bucket: svcCtx.Config.Minio.FormalBucket,
			Object: targetObjectKey,
		}
		if _, err := svcCtx.MinioClient.CopyObject(ctx, destination, source); err != nil {
			return fmt.Errorf("迁移图片素材%d失败: %w", asset.ID, err)
		}
		promoted = append(promoted, promotedImage{
			ProcessAsset:    asset,
			targetObjectKey: targetObjectKey,
		})
	}

	updates := make([]mediarepo.AssetPromotion, 0, len(promoted))
	for _, asset := range promoted {
		updates = append(updates, mediarepo.AssetPromotion{
			AssetID: asset.ID, FormalBucket: svcCtx.Config.Minio.FormalBucket, FormalObjectKey: asset.targetObjectKey,
		})
	}
	if err := svcCtx.MediaRepo.CompleteImagePromotion(ctx, workID, updates); err != nil {
		return fmt.Errorf("更新图片作品处理状态失败: %w", err)
	}

	// 数据库已指向正式对象后，再清理临时桶文件。清理失败不影响作品审核。
	for _, asset := range promoted {
		if asset.Bucket != "" && asset.ObjectKey != "" && asset.Status != "ready" {
			_ = svcCtx.MinioClient.RemoveObject(ctx, asset.Bucket, asset.ObjectKey, minio.RemoveObjectOptions{})
		}
	}
	return nil
}

func markImagePromotionFailed(ctx context.Context, svcCtx *svc.ServiceContext, workID int64, processErr error) {
	message := processErr.Error()
	if len(message) > 1800 {
		message = message[len(message)-1800:]
	}
	_ = svcCtx.MediaRepo.FailProcessing(ctx, workID, message)
}

func ensureFormalMediaBucket(ctx context.Context, svcCtx *svc.ServiceContext) error {
	bucket := svcCtx.Config.Minio.FormalBucket
	if bucket == "" {
		return fmt.Errorf("未配置MinIO正式资源桶")
	}
	exists, err := svcCtx.MinioClient.BucketExists(ctx, bucket)
	if err != nil {
		return fmt.Errorf("检查MinIO正式资源桶失败: %w", err)
	}
	if !exists {
		if err = svcCtx.MinioClient.MakeBucket(ctx, bucket, minio.MakeBucketOptions{}); err != nil {
			return fmt.Errorf("创建MinIO正式资源桶失败: %w", err)
		}
	}
	return nil
}
