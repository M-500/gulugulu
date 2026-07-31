package user_repo

import (
	"context"

	"gorm.io/gorm"
)

type UserDao interface {
	FindOneByEmail(ctx context.Context, email string) (*UserModel, error)
	Insert(ctx context.Context, user *UserModel) error
	Update(ctx context.Context, user *UserModel) error
	Delete(ctx context.Context, user *UserModel) error
}

type userDaoImpl struct {
	db *gorm.DB
}

func NewUserDao(db *gorm.DB) UserDao {
	return &userDaoImpl{
		db: db,
	}
}

func (d *userDaoImpl) FindOneByEmail(ctx context.Context, email string) (*UserModel, error) {
	var user UserModel
	if err := d.db.WithContext(ctx).Where("email = ?", email).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (d *userDaoImpl) Insert(ctx context.Context, user *UserModel) error {
	return d.db.WithContext(ctx).Create(user).Error
}

func (d *userDaoImpl) Update(ctx context.Context, user *UserModel) error {
	return d.db.WithContext(ctx).Save(user).Error
}

func (d *userDaoImpl) Delete(ctx context.Context, user *UserModel) error {
	return d.db.WithContext(ctx).Delete(user).Error
}
