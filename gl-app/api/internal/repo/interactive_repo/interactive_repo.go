package interactive_repo

import (
	"context"
	"gl-app/api/internal/constants"

	"gorm.io/gorm"
)

type InteractiveRepo interface {
	FindOneByResourceID(ctx context.Context, resourceID int64, bizType constants.BizType) (*InteractiveModel, error)
	Like(ctx context.Context, resourceID int64) error
	Dislike(ctx context.Context, resourceID int64) error
	Share(ctx context.Context, resourceID int64) error
	View(ctx context.Context, resourceID int64) error
}

type interactiveRepoImpl struct {
	dao InteractiveDao
}

func NewInteractiveRepo(dao InteractiveDao) InteractiveRepo {
	return &interactiveRepoImpl{dao: dao}
}

func (r *interactiveRepoImpl) FindOneByResourceID(ctx context.Context, resourceID int64, bizType constants.BizType) (*InteractiveModel, error) {
	return r.dao.FindOneByResourceID(ctx, resourceID, bizType)
}

func (r *interactiveRepoImpl) Like(ctx context.Context, resourceID int64) error {
	return r.dao.UpdateByMap(ctx, resourceID, map[string]any{"like_count": gorm.Expr("like_count + ?", 1)})
}

func (r *interactiveRepoImpl) Dislike(ctx context.Context, resourceID int64) error {
	return r.dao.UpdateByMap(ctx, resourceID, map[string]any{"like_count": gorm.Expr("like_count - ?", 1)})
}

func (r *interactiveRepoImpl) Share(ctx context.Context, resourceID int64) error {
	return r.dao.UpdateByMap(ctx, resourceID, map[string]any{"share_count": gorm.Expr("share_count + ?", 1)})
}

func (r *interactiveRepoImpl) View(ctx context.Context, resourceID int64) error {
	return r.dao.UpdateByMap(ctx, resourceID, map[string]any{"view_count": gorm.Expr("view_count + ?", 1)})
}
