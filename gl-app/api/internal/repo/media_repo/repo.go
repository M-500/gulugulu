package media_repo

import (
	"context"
	"fmt"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

var ErrNotFound = gorm.ErrRecordNotFound

// MediaRepo 统一封装素材、处理任务和 Worker 状态流转。
type MediaRepo interface {
	CreateAsset(ctx context.Context, asset *MediaAsset) error
	FindAssetByID(ctx context.Context, assetID int64) (*MediaAsset, error)
	UpdateAsset(ctx context.Context, asset *MediaAsset) error
	DeleteAsset(ctx context.Context, assetID int64) error
	RecoverStaleTasks(ctx context.Context, before time.Time) error
	ClaimTask(ctx context.Context, workID int64) (bool, error)
	StartImagePromotion(ctx context.Context, workID int64) error
	ListProcessAssets(ctx context.Context, workID int64) ([]ProcessAsset, error)
	ListImageAssets(ctx context.Context, workID int64) ([]ProcessAsset, error)
	UpdateTaskProgress(ctx context.Context, workID int64, stage string, progress int64) error
	MarkAssetReady(ctx context.Context, assetID int64, bucket, objectKey string) error
	MarkVideoAssetReady(ctx context.Context, assetID int64, bucket, objectKey string, durationMs, width, height int64) error
	CompleteProcessing(ctx context.Context, workID int64) error
	CompleteImagePromotion(ctx context.Context, workID int64, assets []AssetPromotion) error
	FailProcessing(ctx context.Context, workID int64, message string) error
	ListPendingWorks(ctx context.Context, before time.Time, limit int) ([]PendingWork, error)
	PublishScheduledWorks(ctx context.Context, now time.Time) error
	RetryFailedWork(ctx context.Context, workID, userID int64) (string, error)
}

type mediaRepoImpl struct {
	db *gorm.DB
}

func NewMediaRepo(db *gorm.DB) MediaRepo {
	return &mediaRepoImpl{db: db}
}

func (r *mediaRepoImpl) CreateAsset(ctx context.Context, asset *MediaAsset) error {
	return r.db.WithContext(ctx).Create(asset).Error
}

func (r *mediaRepoImpl) FindAssetByID(ctx context.Context, assetID int64) (*MediaAsset, error) {
	var asset MediaAsset
	err := r.db.WithContext(ctx).First(&asset, assetID).Error
	return &asset, err
}

func (r *mediaRepoImpl) UpdateAsset(ctx context.Context, asset *MediaAsset) error {
	return r.db.WithContext(ctx).Save(asset).Error
}

func (r *mediaRepoImpl) DeleteAsset(ctx context.Context, assetID int64) error {
	return r.db.WithContext(ctx).Delete(&MediaAsset{}, assetID).Error
}

func (r *mediaRepoImpl) RecoverStaleTasks(ctx context.Context, before time.Time) error {
	return r.db.WithContext(ctx).Model(&processTask{}).
		Where("status = ? AND updated_at < ?", "processing", before).
		Updates(map[string]any{"status": "pending", "stage": "queued"}).Error
}

func (r *mediaRepoImpl) ClaimTask(ctx context.Context, workID int64) (bool, error) {
	claimed := false
	err := r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		result := tx.Model(&processTask{}).
			Where("work_id = ? AND status IN ?", workID, []string{"pending", "failed"}).
			Updates(map[string]any{"status": "processing", "stage": "preparing", "progress": 5, "error_message": ""})
		if result.Error != nil || result.RowsAffected == 0 {
			return result.Error
		}
		claimed = true
		return tx.Model(&workState{}).Where("id = ?", workID).
			Updates(map[string]any{"process_status": "processing", "review_status": "waiting_process"}).Error
	})
	return claimed, err
}

func (r *mediaRepoImpl) StartImagePromotion(ctx context.Context, workID int64) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&workState{}).Where("id = ?", workID).
			Updates(map[string]any{"process_status": "processing", "review_status": "waiting_process"}).Error; err != nil {
			return err
		}
		return tx.Model(&processTask{}).Where("work_id = ?", workID).
			Updates(map[string]any{"status": "processing", "stage": "promoting_images", "progress": 10, "error_message": ""}).Error
	})
}

func (r *mediaRepoImpl) ListProcessAssets(ctx context.Context, workID int64) ([]ProcessAsset, error) {
	return r.listAssets(ctx, workID, nil)
}

func (r *mediaRepoImpl) ListImageAssets(ctx context.Context, workID int64) ([]ProcessAsset, error) {
	return r.listAssets(ctx, workID, []string{"image", "cover"})
}

func (r *mediaRepoImpl) listAssets(ctx context.Context, workID int64, roles []string) ([]ProcessAsset, error) {
	var rows []ProcessAsset
	db := r.db.WithContext(ctx).Table("work_asset AS wa").
		Select("ma.id,wa.role,ma.bucket,ma.object_key,ma.origin_name,ma.content_type,ma.formal_bucket,ma.formal_object_key,ma.status").
		Joins("JOIN media_asset AS ma ON ma.id = wa.media_asset_id AND ma.deleted_at IS NULL").
		Where("wa.work_id = ?", workID)
	if len(roles) > 0 {
		db = db.Where("wa.role IN ?", roles)
	}
	err := db.Order("CASE wa.role WHEN 'cover' THEN 0 ELSE 1 END,wa.sort,wa.id").Scan(&rows).Error
	return rows, err
}

func (r *mediaRepoImpl) UpdateTaskProgress(ctx context.Context, workID int64, stage string, progress int64) error {
	return r.db.WithContext(ctx).Model(&processTask{}).Where("work_id = ?", workID).
		Updates(map[string]any{"stage": stage, "progress": progress}).Error
}

func (r *mediaRepoImpl) MarkAssetReady(ctx context.Context, assetID int64, bucket, objectKey string) error {
	return r.db.WithContext(ctx).Model(&MediaAsset{}).Where("id = ?", assetID).
		Updates(map[string]any{"status": "ready", "formal_bucket": bucket, "formal_object_key": objectKey, "process_error": ""}).Error
}

func (r *mediaRepoImpl) MarkVideoAssetReady(ctx context.Context, assetID int64, bucket, objectKey string, durationMs, width, height int64) error {
	return r.db.WithContext(ctx).Model(&MediaAsset{}).Where("id = ?", assetID).
		Updates(map[string]any{"status": "ready", "formal_bucket": bucket, "formal_object_key": objectKey,
			"duration_ms": durationMs, "width": width, "height": height, "process_error": ""}).Error
}

func (r *mediaRepoImpl) CompleteProcessing(ctx context.Context, workID int64) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&processTask{}).Where("work_id = ?", workID).
			Updates(map[string]any{"status": "succeeded", "stage": "waiting_review", "progress": 100, "error_message": ""}).Error; err != nil {
			return err
		}
		return tx.Model(&workState{}).Where("id = ?", workID).
			Updates(map[string]any{"process_status": "succeeded", "review_status": "pending_review", "publish_status": "pending"}).Error
	})
}

func (r *mediaRepoImpl) CompleteImagePromotion(ctx context.Context, workID int64, assets []AssetPromotion) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		for _, asset := range assets {
			if err := tx.Model(&MediaAsset{}).Where("id = ?", asset.AssetID).
				Updates(map[string]any{"status": "ready", "formal_bucket": asset.FormalBucket,
					"formal_object_key": asset.FormalObjectKey, "process_error": ""}).Error; err != nil {
				return err
			}
		}
		if err := tx.Model(&processTask{}).Where("work_id = ?", workID).
			Updates(map[string]any{"status": "succeeded", "stage": "waiting_review", "progress": 100, "error_message": ""}).Error; err != nil {
			return err
		}
		return tx.Model(&workState{}).Where("id = ?", workID).
			Updates(map[string]any{"process_status": "succeeded", "review_status": "pending_review", "publish_status": "pending"}).Error
	})
}

func (r *mediaRepoImpl) FailProcessing(ctx context.Context, workID int64, message string) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&processTask{}).Where("work_id = ?", workID).
			Updates(map[string]any{"status": "failed", "stage": "failed", "error_message": message}).Error; err != nil {
			return err
		}
		if err := tx.Model(&workState{}).Where("id = ?", workID).
			Updates(map[string]any{"process_status": "failed", "review_status": "waiting_process", "publish_status": "pending"}).Error; err != nil {
			return err
		}
		subQuery := tx.Table("work_asset").Select("media_asset_id").Where("work_id = ?", workID)
		return tx.Model(&MediaAsset{}).Where("id IN (?) AND status <> ?", subQuery, "ready").
			Updates(map[string]any{"status": "failed", "process_error": message}).Error
	})
}

func (r *mediaRepoImpl) ListPendingWorks(ctx context.Context, before time.Time, limit int) ([]PendingWork, error) {
	var rows []PendingWork
	err := r.db.WithContext(ctx).Table("media_process_task AS t").Select("t.work_id,w.type").
		Joins("JOIN work AS w ON w.id = t.work_id AND w.deleted_at IS NULL").
		Where("t.status = ? AND t.updated_at < ?", "pending", before).Limit(limit).Scan(&rows).Error
	return rows, err
}

func (r *mediaRepoImpl) PublishScheduledWorks(ctx context.Context, now time.Time) error {
	return r.db.WithContext(ctx).Model(&workState{}).
		Where("review_status = ? AND publish_status = ? AND scheduled_at <= ?", "approved", "scheduled", now).
		Updates(map[string]any{"publish_status": "published", "published_at": now}).Error
}

func (r *mediaRepoImpl) RetryFailedWork(ctx context.Context, workID, userID int64) (string, error) {
	var workType string
	err := r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var work workState
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
			Where("id = ? AND user_id = ? AND process_status = ?", workID, userID, "failed").First(&work).Error; err != nil {
			return err
		}
		workType = work.Type
		if err := tx.Model(&workState{}).Where("id = ?", workID).
			Updates(map[string]any{"process_status": "pending", "review_status": "waiting_process"}).Error; err != nil {
			return err
		}
		result := tx.Model(&processTask{}).Where("work_id = ?", workID).
			Updates(map[string]any{"status": "pending", "stage": "queued", "progress": 0,
				"error_message": "", "retry_count": gorm.Expr("retry_count + 1")})
		if result.Error != nil {
			return result.Error
		}
		if result.RowsAffected == 0 {
			return fmt.Errorf("媒体处理任务不存在")
		}
		return nil
	})
	return workType, err
}
