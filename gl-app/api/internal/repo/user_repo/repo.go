package user_repo

import (
	"context"
	"errors"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/zeromicro/go-zero/core/logx"
)

const userCacheTTL = 12 * time.Hour

// UserRepo 是业务层唯一依赖的用户持久化接口。
type UserRepo interface {
	FindOneByEmail(ctx context.Context, email string) (*User, error)
	FindOneByID(ctx context.Context, userID int64) (*User, error)
	Create(ctx context.Context, user *User) error
	Update(ctx context.Context, user *User) error
	UpdateByMap(ctx context.Context, userID int64, data map[string]any) error
	Delete(ctx context.Context, userID int64) error
}

type userRepoImpl struct {
	dao   UserDAO
	cache UserCache
}

func NewUserRepo(dao UserDAO, cache UserCache) UserRepo {
	return &userRepoImpl{dao: dao, cache: cache}
}

func (u *userRepoImpl) FindOneByID(ctx context.Context, userID int64) (*User, error) {
	user, err := u.cache.FindByID(ctx, userID)
	if err == nil {
		return user, nil
	}
	if !errors.Is(err, redis.Nil) {
		// Redis 故障不应阻断主业务，降级读取数据库并记录错误。
		logx.WithContext(ctx).Errorf("读取用户缓存失败，已降级查询数据库: %v", err)
	}

	user, err = u.dao.FindOneByID(ctx, userID)
	if err != nil {
		return nil, err
	}
	if cacheErr := u.cache.SetByID(ctx, userID, user, userCacheTTL); cacheErr != nil {
		logx.WithContext(ctx).Errorf("回写用户缓存失败: %v", cacheErr)
	}
	return user, nil
}

func (u *userRepoImpl) FindOneByEmail(ctx context.Context, email string) (*User, error) {
	return u.dao.FindOneByEmail(ctx, email)
}

func (u *userRepoImpl) Create(ctx context.Context, user *User) error {
	if err := u.dao.Create(ctx, user); err != nil {
		return err
	}
	u.deleteCache(ctx, user.ID)
	return nil
}

func (u *userRepoImpl) Update(ctx context.Context, user *User) error {
	if err := u.dao.Update(ctx, user); err != nil {
		return err
	}
	u.deleteCache(ctx, user.ID)
	return nil
}

func (u *userRepoImpl) UpdateByMap(ctx context.Context, userID int64, data map[string]any) error {
	if err := u.dao.UpdateByMap(ctx, userID, data); err != nil {
		return err
	}
	u.deleteCache(ctx, userID)
	return nil
}

func (u *userRepoImpl) Delete(ctx context.Context, userID int64) error {
	if err := u.dao.Delete(ctx, userID); err != nil {
		return err
	}
	u.deleteCache(ctx, userID)
	return nil
}

func (u *userRepoImpl) deleteCache(ctx context.Context, userID int64) {
	if err := u.cache.DeleteByID(ctx, userID); err != nil {
		logx.WithContext(ctx).Errorf("删除用户缓存失败: %v", err)
	}
}
