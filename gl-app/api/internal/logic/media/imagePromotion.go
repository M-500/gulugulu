package media

import (
	"context"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/minio/minio-go/v7"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"gl-app/api/internal/svc"
)

type imagePromotionAsset struct {
	ID              int64  `db:"id"`
	Role            string `db:"role"`
	Bucket          string `db:"bucket"`
	ObjectKey       string `db:"object_key"`
	OriginName      string `db:"origin_name"`
	Status          string `db:"status"`
	FormalObjectKey string `db:"formal_object_key"`
}

type promotedImage struct {
	imagePromotionAsset
	targetObjectKey string
}

// promoteImageWork 在MinIO服务端直接复制图片，不需要下载文件，也不经过Kafka。
// 全部图片和封面复制成功并写入数据库后，作品立即进入待审核状态。
func promoteImageWork(ctx context.Context, svcCtx *svc.ServiceContext, workID int64) error {
	if err := ensureFormalMediaBucket(ctx, svcCtx); err != nil {
		return err
	}
	_, _ = svcCtx.SqlConn.ExecCtx(ctx, `UPDATE work SET process_status='processing',
		review_status='waiting_process' WHERE id=?`, workID)
	_, _ = svcCtx.SqlConn.ExecCtx(ctx, `UPDATE media_process_task SET status='processing',
		stage='promoting_images',progress=10,error_message='' WHERE work_id=?`, workID)

	var assets []imagePromotionAsset
	if err := svcCtx.SqlConn.QueryRowsCtx(ctx, &assets, `SELECT ma.id,wa.role,ma.bucket,ma.object_key,
		ma.origin_name,ma.status,ma.formal_object_key
		FROM work_asset wa
		JOIN media_asset ma ON ma.id=wa.media_asset_id
		WHERE wa.work_id=? AND wa.role IN ('image','cover')
		ORDER BY CASE wa.role WHEN 'cover' THEN 0 ELSE 1 END,wa.sort`, workID); err != nil {
		return fmt.Errorf("查询图片作品素材失败: %w", err)
	}
	if len(assets) < 2 {
		return fmt.Errorf("图片作品素材不完整")
	}

	promoted := make([]promotedImage, 0, len(assets))
	for _, asset := range assets {
		if asset.Status == "ready" && asset.FormalObjectKey != "" {
			promoted = append(promoted, promotedImage{
				imagePromotionAsset: asset,
				targetObjectKey:     asset.FormalObjectKey,
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
			imagePromotionAsset: asset,
			targetObjectKey:     targetObjectKey,
		})
	}

	if err := svcCtx.SqlConn.TransactCtx(ctx, func(ctx context.Context, session sqlx.Session) error {
		for _, asset := range promoted {
			if _, err := session.ExecCtx(ctx, `UPDATE media_asset SET status='ready',formal_bucket=?,
				formal_object_key=?,process_error='' WHERE id=?`,
				svcCtx.Config.Minio.FormalBucket, asset.targetObjectKey, asset.ID); err != nil {
				return err
			}
		}
		if _, err := session.ExecCtx(ctx, `UPDATE media_process_task SET status='succeeded',
			stage='waiting_review',progress=100,error_message='' WHERE work_id=?`, workID); err != nil {
			return err
		}
		_, err := session.ExecCtx(ctx, `UPDATE work SET process_status='succeeded',
			review_status='pending_review',publish_status='pending' WHERE id=?`, workID)
		return err
	}); err != nil {
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
	_, _ = svcCtx.SqlConn.ExecCtx(ctx, `UPDATE media_process_task SET status='failed',
		stage='failed',error_message=? WHERE work_id=?`, message, workID)
	_, _ = svcCtx.SqlConn.ExecCtx(ctx, `UPDATE work SET process_status='failed',
		review_status='waiting_process',publish_status='pending' WHERE id=?`, workID)
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
