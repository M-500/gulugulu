package interactive_repo

import (
	"context"
	"errors"

	"gl-app/api/internal/constants"

	"gorm.io/gorm"
)

var ErrResourceNotFound = errors.New("interactive resource not found")

type InteractiveRepo interface {
	FindOneByResourceID(ctx context.Context, resourceID int64, bizType constants.BizType) (*InteractiveModel, error)
	FindByResourceIDs(ctx context.Context, resourceIDs []int64, bizType constants.BizType) (map[int64]*InteractiveModel, error)
	FindUserLikedResourceIDs(ctx context.Context, userID int64, resourceIDs []int64, bizType constants.BizType) (map[int64]bool, error)
	ResourceExists(ctx context.Context, resourceID int64, bizType constants.BizType) (bool, error)
	MutateLike(ctx context.Context, userID, resourceID int64, bizType constants.BizType, liked bool) (LikeMutation, error)
	ApplyLikeEvent(ctx context.Context, event LikeEvent) error
	Share(ctx context.Context, resourceID int64, bizType constants.BizType) error
	View(ctx context.Context, resourceID int64, bizType constants.BizType) error
}

type interactiveRepoImpl struct {
	dao   InteractiveDAO
	cache LikeStateCache
}

func NewInteractiveRepo(dao InteractiveDAO, cache LikeStateCache) InteractiveRepo {
	return &interactiveRepoImpl{dao: dao, cache: cache}
}

func (r *interactiveRepoImpl) FindOneByResourceID(ctx context.Context, resourceID int64, bizType constants.BizType) (*InteractiveModel, error) {
	return r.dao.FindOneByResourceID(ctx, resourceID, bizType)
}

func (r *interactiveRepoImpl) FindByResourceIDs(ctx context.Context, resourceIDs []int64, bizType constants.BizType) (map[int64]*InteractiveModel, error) {
	return r.dao.FindByResourceIDs(ctx, resourceIDs, bizType)
}

func (r *interactiveRepoImpl) FindUserLikedResourceIDs(ctx context.Context, userID int64, resourceIDs []int64, bizType constants.BizType) (map[int64]bool, error) {
	states, err := r.dao.FindUserLikedResourceIDs(ctx, userID, resourceIDs, bizType)
	if err != nil { return nil, err }
	latest, cacheErr := r.cache.OverlayUserLikedStates(ctx, userID, resourceIDs, bizType, states)
	if cacheErr != nil { return states, nil }
	return latest, nil
}

func (r *interactiveRepoImpl) ResourceExists(ctx context.Context, resourceID int64, bizType constants.BizType) (bool, error) {
	return r.dao.ResourceExists(ctx, resourceID, bizType)
}

func (r *interactiveRepoImpl) MutateLike(ctx context.Context, userID, resourceID int64, bizType constants.BizType, liked bool) (LikeMutation, error) {
	mutation, hit, err := r.cache.MutateIfPresent(ctx, userID, resourceID, bizType, liked)
	if err != nil {
		return LikeMutation{}, err
	}
	if hit {
		return mutation, nil
	}

	// 只有缓存冷启动时读取数据库；后续高频操作完全在 Redis Lua 中原子完成。
	exists, err := r.dao.ResourceExists(ctx, resourceID, bizType)
	if err != nil {
		return LikeMutation{}, err
	}
	if !exists {
		return LikeMutation{}, ErrResourceNotFound
	}
	initialLiked, initialVersion, err := r.dao.FindUserLike(ctx, userID, resourceID, bizType)
	if err != nil {
		return LikeMutation{}, err
	}
	initialCount := int64(0)
	row, err := r.dao.FindOneByResourceID(ctx, resourceID, bizType)
	if err == nil {
		initialCount = row.LikeCount
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		return LikeMutation{}, err
	}
	return r.cache.Mutate(ctx, userID, resourceID, bizType, liked, initialLiked, initialCount, initialVersion)
}

func (r *interactiveRepoImpl) ApplyLikeEvent(ctx context.Context, event LikeEvent) error {
	return r.dao.ApplyLikeEvent(ctx, event)
}

func (r *interactiveRepoImpl) Share(ctx context.Context, resourceID int64, bizType constants.BizType) error {
	return r.dao.UpdateCounter(ctx, resourceID, bizType, "share_count", 1)
}

func (r *interactiveRepoImpl) View(ctx context.Context, resourceID int64, bizType constants.BizType) error {
	return r.dao.UpdateCounter(ctx, resourceID, bizType, "view_count", 1)
}
