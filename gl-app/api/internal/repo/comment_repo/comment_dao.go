package comment_repo

import (
	"context"
	"errors"

	"gl-app/api/internal/constants"
	interactiverepo "gl-app/api/internal/repo/interactive_repo"
	mediarepo "gl-app/api/internal/repo/media_repo"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

var (
	ErrWorkNotFound     = errors.New("published work not found")
	ErrParentNotFound   = errors.New("parent comment not found")
	ErrImageUnavailable = errors.New("comment image unavailable")
)

type CommentDAO interface {
	Create(ctx context.Context, input CreateInput) (*CommentModel, error)
	FindByID(ctx context.Context, commentID int64) (*CommentView, error)
	ListRoots(ctx context.Context, workID, limit, offset int64) (int64, []CommentView, error)
	ListReplies(ctx context.Context, rootID, limit, offset int64) (int64, []CommentView, error)
}

type commentDAOImpl struct{ db *gorm.DB }

func NewCommentDAO(db *gorm.DB) CommentDAO { return &commentDAOImpl{db: db} }

func (d *commentDAOImpl) Create(ctx context.Context, input CreateInput) (*CommentModel, error) {
	comment := &CommentModel{WorkID: input.WorkID, UserID: input.UserID, ParentID: input.ParentID, Content: input.Content, ImageMediaID: input.ImageMediaID}
	err := d.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var workCount int64
		if err := tx.Table("work").Where("id = ? AND deleted_at IS NULL AND publish_status = ?", input.WorkID, "published").Count(&workCount).Error; err != nil {
			return err
		}
		if workCount == 0 {
			return ErrWorkNotFound
		}

		if input.ParentID > 0 {
			var parent CommentModel
			if err := tx.Clauses(clause.Locking{Strength: "SHARE"}).Where("id = ? AND work_id = ?", input.ParentID, input.WorkID).First(&parent).Error; err != nil {
				if errors.Is(err, gorm.ErrRecordNotFound) {
					return ErrParentNotFound
				}
				return err
			}
			comment.ReplyToUserID = parent.UserID
			if parent.ParentID == 0 {
				comment.RootID = int64(parent.ID)
			} else {
				comment.RootID = parent.RootID
			}
		}

		if err := tx.Create(comment).Error; err != nil {
			return err
		}
		if input.ImageMediaID > 0 {
			result := tx.Model(&mediarepo.MediaAsset{}).
				Where("id = ? AND user_id = ? AND resource_type = ? AND status = ? AND bound_comment_id = 0", input.ImageMediaID, input.UserID, "image", "uploaded").
				Updates(map[string]any{"status": "ready", "formal_bucket": input.ImageBucket, "formal_object_key": input.ImageObjectKey, "bound_comment_id": comment.ID})
			if result.Error != nil {
				return result.Error
			}
			if result.RowsAffected == 0 {
				return ErrImageUnavailable
			}
		}

		interactive := interactiverepo.InteractiveModel{ResourceID: input.WorkID, ResourceType: constants.WorkType, CommentCount: 1}
		return tx.Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "resource_id"}, {Name: "resource_type"}},
			DoUpdates: clause.Assignments(map[string]any{"comment_count": gorm.Expr("comment_count + 1")}),
		}).Create(&interactive).Error
	})
	return comment, err
}

func (d *commentDAOImpl) FindByID(ctx context.Context, commentID int64) (*CommentView, error) {
	var row CommentView
	err := commentViewQuery(d.db.WithContext(ctx).Table("comment AS c").Where("c.id = ? AND c.deleted_at IS NULL", commentID)).
		Select(commentViewSelect() + ", (SELECT COUNT(1) FROM comment r WHERE r.root_id = c.id AND r.deleted_at IS NULL) AS reply_count").
		Take(&row).Error
	return &row, err
}

func (d *commentDAOImpl) ListRoots(ctx context.Context, workID, limit, offset int64) (int64, []CommentView, error) {
	base := d.db.WithContext(ctx).Table("comment AS c").Where("c.work_id = ? AND c.parent_id = 0 AND c.deleted_at IS NULL", workID)
	var total int64
	if err := base.Count(&total).Error; err != nil {
		return 0, nil, err
	}
	var rows []CommentView
	err := commentViewQuery(base).
		Select(commentViewSelect() + ", (SELECT COUNT(1) FROM comment r WHERE r.root_id = c.id AND r.deleted_at IS NULL) AS reply_count").
		Order("c.created_at DESC, c.id DESC").Limit(int(limit)).Offset(int(offset)).Scan(&rows).Error
	return total, rows, err
}

func (d *commentDAOImpl) ListReplies(ctx context.Context, rootID, limit, offset int64) (int64, []CommentView, error) {
	base := d.db.WithContext(ctx).Table("comment AS c").Where("c.root_id = ? AND c.deleted_at IS NULL", rootID)
	var total int64
	if err := base.Count(&total).Error; err != nil {
		return 0, nil, err
	}
	var rows []CommentView
	err := commentViewQuery(base).Select(commentViewSelect() + ", 0 AS reply_count").
		Order("c.created_at ASC, c.id ASC").Limit(int(limit)).Offset(int(offset)).Scan(&rows).Error
	return total, rows, err
}

func commentViewQuery(db *gorm.DB) *gorm.DB {
	return db.Joins("JOIN user AS u ON u.id = c.user_id AND u.deleted_at IS NULL").
		Joins("LEFT JOIN user AS ru ON ru.id = c.reply_to_user_id AND ru.deleted_at IS NULL").
		Joins("LEFT JOIN media_asset AS ma ON ma.id = c.image_media_id AND ma.deleted_at IS NULL").
		Joins("LEFT JOIN interactive AS i ON i.resource_id = c.id AND i.resource_type = ? AND i.deleted_at IS NULL", constants.CommentType)
}

func commentViewSelect() string {
	return "c.id,c.work_id,c.user_id,c.root_id,c.parent_id,c.reply_to_user_id,c.content,c.image_media_id,c.created_at," +
		"u.nickname AS author_name,u.avatar AS author_avatar,COALESCE(ru.nickname,'') AS reply_to_name," +
		"COALESCE(ma.formal_bucket,'') AS image_bucket,COALESCE(ma.formal_object_key,'') AS image_object_key,COALESCE(i.like_count,0) AS like_count"
}
