package interactive_repo

import (
	"context"
	"gl-app/api/internal/constants"

	"gorm.io/gorm"
)

type InteractiveDao interface {
	FindOneByResourceID(ctx context.Context, resourceID int64, bizType constants.BizType) (*InteractiveModel, error)
	Create(ctx context.Context, interactive *InteractiveModel) error
	Update(ctx context.Context, interactive *InteractiveModel) error
	UpdateByMap(ctx context.Context, resourceID int64, data map[string]any) error
	Delete(ctx context.Context, resourceID int64) error
}

type interactiveDaoImpl struct {
	db *gorm.DB
}

func NewInteractiveDao(db *gorm.DB) InteractiveDao {
	return &interactiveDaoImpl{db: db}
}

func (d *interactiveDaoImpl) FindOneByResourceID(ctx context.Context, resourceID int64, bizType constants.BizType) (*InteractiveModel, error) {
	var interactive InteractiveModel
	if err := d.db.WithContext(ctx).Where("resource_id = ? AND resource_type = ?", resourceID, bizType).First(&interactive).Error; err != nil {
		return nil, err
	}
	return &interactive, nil
}

func (d *interactiveDaoImpl) Create(ctx context.Context, interactive *InteractiveModel) error {
	return d.db.WithContext(ctx).Create(interactive).Error
}

func (d *interactiveDaoImpl) Update(ctx context.Context, interactive *InteractiveModel) error {
	return d.db.WithContext(ctx).Save(interactive).Error
}

func (d *interactiveDaoImpl) UpdateByMap(ctx context.Context, resourceID int64, data map[string]any) error {
	result := d.db.WithContext(ctx).Model(&InteractiveModel{}).Where("resource_id = ?", resourceID).Updates(data)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

func (d *interactiveDaoImpl) Delete(ctx context.Context, resourceID int64) error {
	result := d.db.WithContext(ctx).Where("resource_id = ?", resourceID).Delete(&InteractiveModel{})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}
