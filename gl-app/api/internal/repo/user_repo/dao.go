package user_repo

import (
	"context"

	"gorm.io/gorm"
)

// UserDAO 只负责数据库读写，缓存一致性由上层 UserRepo 统一处理。
type UserDAO interface {
	FindOneByEmail(ctx context.Context, email string) (*User, error)
	FindOneByID(ctx context.Context, userID int64) (*User, error)
	Create(ctx context.Context, user *User) error
	Update(ctx context.Context, user *User) error
	UpdateByMap(ctx context.Context, userID int64, data map[string]any) error
	Delete(ctx context.Context, userID int64) error
}

type userDAOImpl struct {
	db *gorm.DB
}

func NewUserDAO(db *gorm.DB) UserDAO {
	return &userDAOImpl{db: db}
}

func (d *userDAOImpl) FindOneByID(ctx context.Context, userID int64) (*User, error) {
	var user User
	if err := d.db.WithContext(ctx).First(&user, userID).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (d *userDAOImpl) FindOneByEmail(ctx context.Context, email string) (*User, error) {
	var user User
	if err := d.db.WithContext(ctx).Where("email = ?", email).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (d *userDAOImpl) Create(ctx context.Context, user *User) error {
	return d.db.WithContext(ctx).Create(user).Error
}

func (d *userDAOImpl) Update(ctx context.Context, user *User) error {
	return d.db.WithContext(ctx).Save(user).Error
}

func (d *userDAOImpl) UpdateByMap(ctx context.Context, userID int64, data map[string]any) error {
	result := d.db.WithContext(ctx).Model(&User{}).Where("id = ?", userID).Updates(data)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return ErrNotFound
	}
	return nil
}

func (d *userDAOImpl) Delete(ctx context.Context, userID int64) error {
	result := d.db.WithContext(ctx).Delete(&User{}, userID)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return ErrNotFound
	}
	return nil
}
