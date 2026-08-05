package interactive_repo

import (
	"context"
	"errors"

	"gl-app/api/internal/constants"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// InteractiveDAO 只负责数据库访问，业务逻辑和缓存编排由 Repo 统一处理。
type InteractiveDAO interface {
	FindOneByResourceID(ctx context.Context, resourceID int64, bizType constants.BizType) (*InteractiveModel, error)
	FindByResourceIDs(ctx context.Context, resourceIDs []int64, bizType constants.BizType) (map[int64]*InteractiveModel, error)
	FindUserLike(ctx context.Context, userID, resourceID int64, bizType constants.BizType) (bool, int64, error)
	FindUserLikedResourceIDs(ctx context.Context, userID int64, resourceIDs []int64, bizType constants.BizType) (map[int64]bool, error)
	ResourceExists(ctx context.Context, resourceID int64, bizType constants.BizType) (bool, error)
	ApplyLikeEvent(ctx context.Context, event LikeEvent) error
	UpdateCounter(ctx context.Context, resourceID int64, bizType constants.BizType, field string, delta int64) error
}

type interactiveDAOImpl struct{ db *gorm.DB }

func NewInteractiveDAO(db *gorm.DB) InteractiveDAO { return &interactiveDAOImpl{db: db} }

func (d *interactiveDAOImpl) FindOneByResourceID(ctx context.Context, resourceID int64, bizType constants.BizType) (*InteractiveModel, error) {
	var row InteractiveModel
	err := d.db.WithContext(ctx).Where("resource_id = ? AND resource_type = ?", resourceID, bizType).First(&row).Error
	return &row, err
}

func (d *interactiveDAOImpl) FindByResourceIDs(ctx context.Context, resourceIDs []int64, bizType constants.BizType) (map[int64]*InteractiveModel, error) {
	result := make(map[int64]*InteractiveModel, len(resourceIDs))
	if len(resourceIDs) == 0 {
		return result, nil
	}
	var rows []*InteractiveModel
	if err := d.db.WithContext(ctx).Where("resource_id IN ? AND resource_type = ?", resourceIDs, bizType).Find(&rows).Error; err != nil {
		return nil, err
	}
	for _, row := range rows {
		result[row.ResourceID] = row
	}
	return result, nil
}

func (d *interactiveDAOImpl) FindUserLike(ctx context.Context, userID, resourceID int64, bizType constants.BizType) (bool, int64, error) {
	var row UserLikeModel
	err := d.db.WithContext(ctx).Where("user_id = ? AND resource_id = ? AND resource_type = ?", userID, resourceID, bizType).First(&row).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return false, 0, nil
	}
	return row.Liked, row.Version, err
}

func (d *interactiveDAOImpl) FindUserLikedResourceIDs(ctx context.Context, userID int64, resourceIDs []int64, bizType constants.BizType) (map[int64]bool, error) {
	result := make(map[int64]bool, len(resourceIDs))
	if userID <= 0 || len(resourceIDs) == 0 {
		return result, nil
	}
	var rows []UserLikeModel
	if err := d.db.WithContext(ctx).Select("resource_id").
		Where("user_id = ? AND resource_type = ? AND resource_id IN ? AND liked = ?", userID, bizType, resourceIDs, true).
		Find(&rows).Error; err != nil {
		return nil, err
	}
	for _, row := range rows {
		result[row.ResourceID] = true
	}
	return result, nil
}

func (d *interactiveDAOImpl) ResourceExists(ctx context.Context, resourceID int64, bizType constants.BizType) (bool, error) {
	var count int64
	db := d.db.WithContext(ctx)
	switch bizType {
	case constants.WorkType:
		db = db.Table("work").Where("id = ? AND deleted_at IS NULL AND publish_status = ?", resourceID, "published")
	case constants.CommentType:
		db = db.Table("comment").Where("id = ? AND deleted_at IS NULL", resourceID)
	case constants.AvatarType:
		db = db.Table("user").Where("id = ? AND deleted_at IS NULL", resourceID)
	default:
		return false, nil
	}
	err := db.Count(&count).Error
	return count > 0, err
}

func (d *interactiveDAOImpl) ApplyLikeEvent(ctx context.Context, event LikeEvent) error {
	return d.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var relation UserLikeModel
		err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
			Where("user_id = ? AND resource_id = ? AND resource_type = ?", event.UserID, event.ResourceID, event.ResourceType).
			First(&relation).Error
		delta := int64(0)
		switch {
		case errors.Is(err, gorm.ErrRecordNotFound):
			relation = UserLikeModel{UserID: event.UserID, ResourceID: event.ResourceID, ResourceType: event.ResourceType, Liked: event.Liked, Version: event.Version}
			if err = tx.Create(&relation).Error; err != nil {
				return err
			}
			if event.Liked {
				delta = 1
			}
		case err != nil:
			return err
		case relation.Version >= event.Version:
			return nil
		default:
			if relation.Liked != event.Liked {
				if event.Liked {
					delta = 1
				} else {
					delta = -1
				}
			}
			if err = tx.Model(&relation).Updates(map[string]any{"liked": event.Liked, "version": event.Version}).Error; err != nil {
				return err
			}
		}

		initial := delta
		if initial < 0 {
			initial = 0
		}
		row := InteractiveModel{ResourceID: event.ResourceID, ResourceType: event.ResourceType, LikeCount: initial}
		return tx.Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "resource_id"}, {Name: "resource_type"}},
			DoUpdates: clause.Assignments(map[string]any{"like_count": gorm.Expr("GREATEST(like_count + ?, 0)", delta)}),
		}).Create(&row).Error
	})
}

func (d *interactiveDAOImpl) UpdateCounter(ctx context.Context, resourceID int64, bizType constants.BizType, field string, delta int64) error {
	row := InteractiveModel{ResourceID: resourceID, ResourceType: bizType}
	switch field {
	case "share_count":
		row.ShareCount = max(delta, 0)
	case "view_count":
		row.ViewCount = max(delta, 0)
	case "comment_count":
		row.CommentCount = max(delta, 0)
	default:
		return errors.New("unsupported interactive counter")
	}
	return d.db.WithContext(ctx).Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "resource_id"}, {Name: "resource_type"}},
		DoUpdates: clause.Assignments(map[string]any{field: gorm.Expr("GREATEST("+field+" + ?, 0)", delta)}),
	}).Create(&row).Error
}
