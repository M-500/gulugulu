package work_repo

import (
	"context"
	"fmt"
	"strings"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

var ErrNotFound = gorm.ErrRecordNotFound

// WorkRepo 封装作品领域的全部数据库访问，业务层只面向此接口编程。
type WorkRepo interface {
	FindByIdempotencyKey(ctx context.Context, userID int64, key string) (*Work, error)
	Create(ctx context.Context, input CreateWorkInput) (int64, error)
	FindOwnerDetail(ctx context.Context, workID, userID int64) (*OwnerWorkDetail, error)
	ListReadyAssets(ctx context.Context, workID int64) ([]AssetView, error)
	ListTopics(ctx context.Context, workID int64) ([]TopicView, error)
	ListCreatorWorks(ctx context.Context, query CreatorListQuery) (WorkSummary, int64, []CreatorWorkItem, error)
	FindWorkStatus(ctx context.Context, workID, userID int64) (*WorkStatus, error)
	UpdateCreatorTitle(ctx context.Context, workID, userID int64, title string) (bool, error)
	UpdateCreatorVisibility(ctx context.Context, workID, userID int64, visibility string) (bool, error)
	DeleteCreatorWork(ctx context.Context, workID, userID int64) (bool, error)
	FindVideoAsset(ctx context.Context, workID, userID int64, scope VideoAccessScope) (*VideoAsset, error)
	ListAuditWorks(ctx context.Context, query AuditListQuery) (AuditSummary, int64, []AuditWorkItem, error)
	FindAuditDetail(ctx context.Context, workID int64) (*AuditWorkDetail, error)
	FindAuditState(ctx context.Context, workID int64) (*AuditState, error)
	Review(ctx context.Context, input ReviewInput) (bool, error)
	ListRecommend(ctx context.Context, limit, offset int64) ([]PublicWorkItem, error)
	ListPublishedByUser(ctx context.Context, userID, limit, offset int64) (int64, []PublicWorkItem, error)
	FindPublicDetail(ctx context.Context, workID int64) (*PublicWorkDetail, error)
}

type workRepoImpl struct {
	db *gorm.DB
}

func NewWorkRepo(db *gorm.DB) WorkRepo {
	return &workRepoImpl{db: db}
}

func (r *workRepoImpl) FindByIdempotencyKey(ctx context.Context, userID int64, key string) (*Work, error) {
	var work Work
	err := r.db.WithContext(ctx).Where("user_id = ? AND idempotency_key = ?", userID, key).First(&work).Error
	return &work, err
}

func (r *workRepoImpl) Create(ctx context.Context, input CreateWorkInput) (int64, error) {
	var workID int64
	err := r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if input.CollectionID > 0 {
			var count int64
			if err := tx.Model(&collection{}).Where("id = ? AND user_id = ?", input.CollectionID, input.UserID).Count(&count).Error; err != nil {
				return err
			}
			if count == 0 {
				return fmt.Errorf("合集不存在或不属于当前用户")
			}
		}

		for _, item := range input.Assets {
			var asset mediaAsset
			if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).First(&asset, item.MediaID).Error; err != nil {
				return fmt.Errorf("素材%d不存在: %w", item.MediaID, err)
			}
			if asset.UserID != input.UserID {
				return fmt.Errorf("素材%d不属于当前用户", item.MediaID)
			}
			if asset.Status != "uploaded" || asset.BoundWorkID != 0 {
				return fmt.Errorf("素材%d状态不可发布", item.MediaID)
			}
			if asset.ResourceType != input.Type {
				return fmt.Errorf("作品类型与素材%d类型不一致", item.MediaID)
			}
		}

		visibilityUsers := input.VisibilityUserIDs
		work := Work{
			UserID: input.UserID, Type: input.Type, Title: input.Title, Content: input.Content,
			Visibility: input.Visibility, CollectionID: input.CollectionID, Original: input.Original,
			CoverAssetID: input.CoverAssetID, ProcessStatus: "pending", ReviewStatus: "waiting_process",
			PublishStatus: "pending", ScheduledAt: input.ScheduledAt, IdempotencyKey: input.IdempotencyKey,
		}
		if visibilityUsers != "" {
			work.VisibilityUserIDs = &visibilityUsers
		}
		if err := tx.Create(&work).Error; err != nil {
			return fmt.Errorf("创建作品失败: %w", err)
		}
		workID = work.ID

		relations := make([]workAsset, 0, len(input.Assets)+1)
		assetIDs := make([]int64, 0, len(input.Assets)+1)
		for _, item := range input.Assets {
			relations = append(relations, workAsset{WorkID: workID, MediaAssetID: item.MediaID, Role: input.Type, Sort: item.Sort})
			assetIDs = append(assetIDs, item.MediaID)
		}
		relations = append(relations, workAsset{WorkID: workID, MediaAssetID: input.CoverAssetID, Role: "cover", Sort: 0})
		assetIDs = append(assetIDs, input.CoverAssetID)
		if err := tx.Create(&relations).Error; err != nil {
			return err
		}
		if err := tx.Model(&mediaAsset{}).Where("id IN ?", assetIDs).
			Updates(map[string]any{"status": "bound", "bound_work_id": workID}).Error; err != nil {
			return err
		}

		for _, item := range input.Topics {
			name := strings.TrimSpace(strings.TrimPrefix(item.Name, "#"))
			if name == "" {
				continue
			}
			normalized := strings.ToLower(name)
			itemTopic := topic{Name: name, NormalizedName: normalized}
			if err := tx.Clauses(clause.OnConflict{
				Columns:   []clause.Column{{Name: "normalized_name"}},
				DoNothing: true,
			}).Create(&itemTopic).Error; err != nil {
				return err
			}
			// 并发请求可能已经创建同名话题，此时重新读取唯一键对应的ID。
			if itemTopic.ID == 0 {
				if err := tx.Where("normalized_name = ?", normalized).First(&itemTopic).Error; err != nil {
					return err
				}
			}
			relation := workTopic{WorkID: workID, TopicID: itemTopic.ID, Sort: item.Sort}
			if err := tx.Clauses(clause.OnConflict{DoNothing: true}).Create(&relation).Error; err != nil {
				return err
			}
		}

		return tx.Create(&mediaProcessTask{
			WorkID: workID, TaskType: "process_work", Status: "pending", Stage: "queued", Progress: 0,
		}).Error
	})
	return workID, err
}

func (r *workRepoImpl) FindOwnerDetail(ctx context.Context, workID, userID int64) (*OwnerWorkDetail, error) {
	var row OwnerWorkDetail
	result := r.db.WithContext(ctx).Table("work").Where("id = ? AND user_id = ? AND deleted_at IS NULL", workID, userID).Take(&row)
	return &row, recordError(result)
}

func (r *workRepoImpl) ListReadyAssets(ctx context.Context, workID int64) ([]AssetView, error) {
	var rows []AssetView
	err := r.db.WithContext(ctx).Table("work_asset AS wa").
		Select("wa.media_asset_id, wa.role, wa.sort, ma.formal_bucket, ma.formal_object_key, ma.duration_ms, ma.width, ma.height").
		Joins("JOIN media_asset AS ma ON ma.id = wa.media_asset_id").
		Where("wa.work_id = ? AND ma.status = ? AND ma.deleted_at IS NULL", workID, "ready").
		Order("CASE wa.role WHEN 'video' THEN 0 WHEN 'image' THEN 1 ELSE 2 END, wa.sort, wa.id").Scan(&rows).Error
	return rows, err
}

func (r *workRepoImpl) ListTopics(ctx context.Context, workID int64) ([]TopicView, error) {
	var rows []TopicView
	err := r.db.WithContext(ctx).Table("work_topic AS wt").Select("t.id, t.name").
		Joins("JOIN topic AS t ON t.id = wt.topic_id").Where("wt.work_id = ?", workID).
		Order("wt.sort, wt.id").Scan(&rows).Error
	return rows, err
}

func (r *workRepoImpl) ListCreatorWorks(ctx context.Context, query CreatorListQuery) (WorkSummary, int64, []CreatorWorkItem, error) {
	var summary WorkSummary
	err := r.db.WithContext(ctx).Table("work").Select(`COUNT(1) AS all_count,
		COALESCE(SUM(publish_status = 'published'), 0) AS published_count,
		COALESCE(SUM(review_status = 'pending_review'), 0) AS reviewing_count,
		COALESCE(SUM(review_status = 'rejected'), 0) AS rejected_count`).
		Where("user_id = ? AND process_status = ? AND deleted_at IS NULL", query.UserID, "succeeded").Scan(&summary).Error
	if err != nil {
		return summary, 0, nil, err
	}
	base := r.creatorFilter(r.db.WithContext(ctx).Table("work AS w"), query)
	var total int64
	if err = base.Count(&total).Error; err != nil {
		return summary, 0, nil, err
	}
	var rows []CreatorWorkItem
	err = r.creatorFilter(r.db.WithContext(ctx).Table("work AS w"), query).
		Select(`w.id,w.type,w.title,w.visibility,w.review_status,w.publish_status,w.review_reason,
			w.scheduled_at,w.published_at,w.created_at,COALESCE(video_asset.duration_ms,0) AS duration_ms,
			COALESCE(cover.formal_bucket,'') AS cover_bucket,COALESCE(cover.formal_object_key,'') AS cover_object_key`).
		Joins("LEFT JOIN media_asset AS cover ON cover.id = w.cover_asset_id").
		Joins("LEFT JOIN work_asset AS video_relation ON video_relation.work_id = w.id AND video_relation.role = 'video'").
		Joins("LEFT JOIN media_asset AS video_asset ON video_asset.id = video_relation.media_asset_id").
		Order("w.created_at DESC, w.id DESC").Limit(int(query.Limit)).Offset(int(query.Offset)).Scan(&rows).Error
	return summary, total, rows, err
}

func (r *workRepoImpl) creatorFilter(db *gorm.DB, query CreatorListQuery) *gorm.DB {
	db = db.Where("w.user_id = ? AND w.process_status = ? AND w.deleted_at IS NULL", query.UserID, "succeeded")
	switch query.Status {
	case "published":
		db = db.Where("w.publish_status = ?", "published")
	case "reviewing":
		db = db.Where("w.review_status = ?", "pending_review")
	case "rejected":
		db = db.Where("w.review_status = ?", "rejected")
	}
	if query.Keyword != "" {
		db = db.Where("w.title LIKE ?", "%"+escapeLike(query.Keyword)+"%")
	}
	return db
}

func (r *workRepoImpl) FindWorkStatus(ctx context.Context, workID, userID int64) (*WorkStatus, error) {
	var row WorkStatus
	result := r.db.WithContext(ctx).Table("work AS w").
		Select("w.id,w.process_status,w.review_status,w.publish_status,w.scheduled_at,w.published_at,COALESCE(t.stage,'') AS stage,COALESCE(t.progress,0) AS progress,COALESCE(t.error_message,'') AS error_message").
		Joins("LEFT JOIN media_process_task AS t ON t.work_id = w.id AND t.task_type = 'process_work'").
		Where("w.id = ? AND w.user_id = ? AND w.deleted_at IS NULL", workID, userID).Take(&row)
	return &row, recordError(result)
}

func (r *workRepoImpl) UpdateCreatorTitle(ctx context.Context, workID, userID int64, title string) (bool, error) {
	result := r.db.WithContext(ctx).Model(&Work{}).
		Where("id = ? AND user_id = ? AND process_status = ?", workID, userID, "succeeded").Update("title", title)
	return result.RowsAffected == 1, result.Error
}

func (r *workRepoImpl) UpdateCreatorVisibility(ctx context.Context, workID, userID int64, visibility string) (bool, error) {
	result := r.db.WithContext(ctx).Model(&Work{}).
		Where("id = ? AND user_id = ? AND process_status = ?", workID, userID, "succeeded").
		Updates(map[string]any{"visibility": visibility, "visibility_user_ids": nil})
	return result.RowsAffected == 1, result.Error
}

func (r *workRepoImpl) DeleteCreatorWork(ctx context.Context, workID, userID int64) (bool, error) {
	result := r.db.WithContext(ctx).Model(&Work{}).
		Where("id = ? AND user_id = ? AND process_status = ?", workID, userID, "succeeded").
		Updates(map[string]any{"deleted_at": time.Now(), "publish_status": "deleted"})
	return result.RowsAffected == 1, result.Error
}

func (r *workRepoImpl) FindVideoAsset(ctx context.Context, workID, userID int64, scope VideoAccessScope) (*VideoAsset, error) {
	var row VideoAsset
	db := r.db.WithContext(ctx).Table("work AS w").Select("ma.formal_bucket,ma.formal_object_key").
		Joins("JOIN work_asset AS wa ON wa.work_id = w.id AND wa.role = 'video'").
		Joins("JOIN media_asset AS ma ON ma.id = wa.media_asset_id AND ma.deleted_at IS NULL").
		Where("w.id = ? AND w.deleted_at IS NULL AND ma.status = ?", workID, "ready")
	switch scope {
	case VideoAccessOwner:
		db = db.Where("w.user_id = ?", userID)
	case VideoAccessAudit:
		db = db.Where("w.process_status = ?", "succeeded")
	case VideoAccessPublic:
		db = db.Where("w.type = ? AND w.process_status = ? AND w.review_status = ? AND w.publish_status = ? AND w.visibility = ?",
			"video", "succeeded", "approved", "published", "public")
	}
	result := db.Take(&row)
	return &row, recordError(result)
}

func (r *workRepoImpl) ListAuditWorks(ctx context.Context, query AuditListQuery) (AuditSummary, int64, []AuditWorkItem, error) {
	var summary AuditSummary
	err := r.db.WithContext(ctx).Table("work").Select(`COALESCE(SUM(review_status='pending_review'),0) AS pending_count,
		COALESCE(SUM(review_status='approved'),0) AS approved_count,
		COALESCE(SUM(review_status='rejected'),0) AS rejected_count`).
		Where("process_status = ? AND deleted_at IS NULL", "succeeded").Scan(&summary).Error
	if err != nil {
		return summary, 0, nil, err
	}
	base := r.auditFilter(r.db.WithContext(ctx).Table("work AS w").Joins("JOIN user AS u ON u.id = w.user_id AND u.deleted_at IS NULL"), query)
	var total int64
	if err = base.Count(&total).Error; err != nil {
		return summary, 0, nil, err
	}
	var rows []AuditWorkItem
	err = r.auditFilter(r.db.WithContext(ctx).Table("work AS w").Joins("JOIN user AS u ON u.id = w.user_id AND u.deleted_at IS NULL"), query).
		Select(`w.id,w.type,w.title,LEFT(w.content,160) AS content_excerpt,w.visibility,w.user_id,
			COALESCE(NULLIF(u.nickname,''),u.email) AS author_name,w.review_status,w.publish_status,
			w.review_reason,w.created_at,w.reviewed_at,COALESCE(video_asset.duration_ms,0) AS duration_ms,
			COALESCE(cover.formal_bucket,'') AS cover_bucket,COALESCE(cover.formal_object_key,'') AS cover_object_key`).
		Joins("LEFT JOIN media_asset AS cover ON cover.id = w.cover_asset_id").
		Joins("LEFT JOIN work_asset AS video_relation ON video_relation.work_id = w.id AND video_relation.role = 'video'").
		Joins("LEFT JOIN media_asset AS video_asset ON video_asset.id = video_relation.media_asset_id").
		Order("CASE w.review_status WHEN 'pending_review' THEN 0 ELSE 1 END, w.created_at ASC, w.id ASC").
		Limit(int(query.Limit)).Offset(int(query.Offset)).Scan(&rows).Error
	return summary, total, rows, err
}

func (r *workRepoImpl) auditFilter(db *gorm.DB, query AuditListQuery) *gorm.DB {
	db = db.Where("w.process_status = ? AND w.deleted_at IS NULL", "succeeded")
	switch query.Status {
	case "pending":
		db = db.Where("w.review_status = ?", "pending_review")
	case "approved":
		db = db.Where("w.review_status = ?", "approved")
	case "rejected":
		db = db.Where("w.review_status = ?", "rejected")
	default:
		db = db.Where("w.review_status IN ?", []string{"pending_review", "approved", "rejected"})
	}
	if query.WorkType != "" {
		db = db.Where("w.type = ?", query.WorkType)
	}
	if query.Keyword != "" {
		keyword := "%" + escapeLike(query.Keyword) + "%"
		db = db.Where("w.title LIKE ? OR u.nickname LIKE ? OR u.email LIKE ?", keyword, keyword, keyword)
	}
	return db
}

func (r *workRepoImpl) FindAuditDetail(ctx context.Context, workID int64) (*AuditWorkDetail, error) {
	var row AuditWorkDetail
	result := r.db.WithContext(ctx).Table("work AS w").
		Select(`w.id,w.type,w.title,w.content,w.visibility,w.original,w.user_id,
			COALESCE(NULLIF(u.nickname,''),u.email) AS author_name,u.email AS author_email,
			w.review_status,w.publish_status,w.review_reason,w.scheduled_at,w.created_at,w.reviewed_at,
			COALESCE(cover.formal_bucket,'') AS cover_bucket,COALESCE(cover.formal_object_key,'') AS cover_object_key`).
		Joins("JOIN user AS u ON u.id = w.user_id AND u.deleted_at IS NULL").
		Joins("LEFT JOIN media_asset AS cover ON cover.id = w.cover_asset_id").
		Where("w.id = ? AND w.process_status = ? AND w.deleted_at IS NULL", workID, "succeeded").Take(&row)
	return &row, recordError(result)
}

func (r *workRepoImpl) FindAuditState(ctx context.Context, workID int64) (*AuditState, error) {
	var row AuditState
	result := r.db.WithContext(ctx).Table("work").Select("process_status,review_status,scheduled_at").
		Where("id = ? AND deleted_at IS NULL", workID).Take(&row)
	return &row, recordError(result)
}

func (r *workRepoImpl) Review(ctx context.Context, input ReviewInput) (bool, error) {
	now := time.Now().UTC()
	result := r.db.WithContext(ctx).Model(&Work{}).Where("id = ? AND review_status = ?", input.WorkID, "pending_review").
		Updates(map[string]any{"review_status": input.ReviewStatus, "publish_status": input.PublishStatus,
			"reviewed_by": input.ReviewerID, "reviewed_at": now, "review_reason": input.Reason, "published_at": input.PublishedAt})
	return result.RowsAffected == 1, result.Error
}

func (r *workRepoImpl) ListRecommend(ctx context.Context, limit, offset int64) ([]PublicWorkItem, error) {
	var rows []PublicWorkItem
	err := r.publicListBase(ctx).
		Select(`w.id,w.type,w.title,LEFT(w.content,120) AS content_excerpt,w.published_at,
			COALESCE(cover.formal_bucket,'') AS cover_bucket,COALESCE(cover.formal_object_key,'') AS cover_object_key,
			COALESCE(video.duration_ms,0) AS duration_ms,u.id AS author_id,
			COALESCE(NULLIF(u.nickname,''),u.email) AS author_name,u.avatar AS author_avatar`).
		Order("w.published_at DESC,w.id DESC").Limit(int(limit)).Offset(int(offset)).Scan(&rows).Error
	return rows, err
}

func (r *workRepoImpl) ListPublishedByUser(ctx context.Context, userID, limit, offset int64) (int64, []PublicWorkItem, error) {
	base := r.publicListBase(ctx).Where("w.user_id = ?", userID)
	var total int64
	if err := base.Count(&total).Error; err != nil {
		return 0, nil, err
	}
	var rows []PublicWorkItem
	err := r.publicListBase(ctx).Where("w.user_id = ?", userID).
		Select(`w.id,w.type,w.title,LEFT(w.content,120) AS content_excerpt,w.published_at,
			COALESCE(cover.formal_bucket,'') AS cover_bucket,COALESCE(cover.formal_object_key,'') AS cover_object_key,
			COALESCE(video.duration_ms,0) AS duration_ms`).
		Order("w.published_at DESC,w.id DESC").Limit(int(limit)).Offset(int(offset)).Scan(&rows).Error
	return total, rows, err
}

func (r *workRepoImpl) publicListBase(ctx context.Context) *gorm.DB {
	return r.db.WithContext(ctx).Table("work AS w").
		Joins("JOIN user AS u ON u.id = w.user_id AND u.deleted_at IS NULL").
		Joins("LEFT JOIN media_asset AS cover ON cover.id = w.cover_asset_id AND cover.deleted_at IS NULL").
		Joins("LEFT JOIN work_asset AS video_relation ON video_relation.work_id = w.id AND video_relation.role = 'video' AND video_relation.sort = 0").
		Joins("LEFT JOIN media_asset AS video ON video.id = video_relation.media_asset_id AND video.deleted_at IS NULL").
		Where("w.deleted_at IS NULL AND w.process_status = ? AND w.review_status = ? AND w.publish_status = ? AND w.visibility = ?",
			"succeeded", "approved", "published", "public")
}

func (r *workRepoImpl) FindPublicDetail(ctx context.Context, workID int64) (*PublicWorkDetail, error) {
	var row PublicWorkDetail
	result := r.db.WithContext(ctx).Table("work AS w").
		Select(`w.id,w.type,w.title,w.content,w.published_at,u.id AS author_id,
			COALESCE(NULLIF(u.nickname,''),u.email) AS author_name,u.avatar AS author_avatar`).
		Joins("JOIN user AS u ON u.id = w.user_id AND u.deleted_at IS NULL").
		Where("w.id = ? AND w.deleted_at IS NULL AND w.process_status = ? AND w.review_status = ? AND w.publish_status = ? AND w.visibility = ?",
			workID, "succeeded", "approved", "published", "public").Take(&row)
	return &row, recordError(result)
}

func recordError(result *gorm.DB) error {
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return ErrNotFound
	}
	return nil
}

func escapeLike(value string) string {
	replacer := strings.NewReplacer(`\`, `\\`, `%`, `\%`, `_`, `\_`)
	return replacer.Replace(value)
}
