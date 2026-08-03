package user_repo

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

// UserCache 对 Redis 操作做接口抽象，便于单元测试替换为内存实现。
type UserCache interface {
	FindByID(ctx context.Context, userID int64) (*User, error)
	SetByID(ctx context.Context, userID int64, user *User, expiration time.Duration) error
	DeleteByID(ctx context.Context, userID int64) error
}

type userCacheImpl struct {
	cmd redis.Cmdable
}

func NewUserCache(cmd redis.Cmdable) UserCache {
	return &userCacheImpl{cmd: cmd}
}

func (u *userCacheImpl) SetByID(ctx context.Context, userID int64, user *User, expiration time.Duration) error {
	value, err := json.Marshal(user)
	if err != nil {
		return fmt.Errorf("序列化用户缓存失败: %w", err)
	}
	return u.cmd.Set(ctx, u.key(userID), value, expiration).Err()
}

func (u *userCacheImpl) DeleteByID(ctx context.Context, userID int64) error {
	return u.cmd.Del(ctx, u.key(userID)).Err()
}

func (u *userCacheImpl) FindByID(ctx context.Context, userID int64) (*User, error) {
	value, err := u.cmd.Get(ctx, u.key(userID)).Bytes()
	if err != nil {
		return nil, err
	}
	var user User
	if err = json.Unmarshal(value, &user); err != nil {
		return nil, fmt.Errorf("反序列化用户缓存失败: %w", err)
	}
	return &user, nil
}

func (u *userCacheImpl) key(userID int64) string {
	return fmt.Sprintf("gulugulu:user:%d", userID)
}
