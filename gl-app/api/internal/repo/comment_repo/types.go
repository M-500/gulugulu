package comment_repo

import (
	"time"

	"gorm.io/gorm"
)

// CommentModel 只允许两种关系：直接评论作品，或回复某条评论。
// RootID 将任意回复归并到一级评论下，展示层级始终限制为两层。
type CommentModel struct {
	gorm.Model
	WorkID        int64  `gorm:"column:work_id;type:bigint;not null;index:idx_work_root_created,priority:1"`
	UserID        int64  `gorm:"column:user_id;type:bigint;not null;index"`
	RootID        int64  `gorm:"column:root_id;type:bigint;not null;default:0;index:idx_work_root_created,priority:2"`
	ParentID      int64  `gorm:"column:parent_id;type:bigint;not null;default:0;index"`
	ReplyToUserID int64  `gorm:"column:reply_to_user_id;type:bigint;not null;default:0"`
	Content       string `gorm:"column:content;type:varchar(1000);not null;default:''"`
	ImageMediaID  int64  `gorm:"column:image_media_id;type:bigint;not null;default:0"`
}

func (CommentModel) TableName() string { return "comment" }

type CreateInput struct {
	WorkID, UserID, ParentID, ImageMediaID int64
	Content                                string
	ImageBucket, ImageObjectKey            string
}

type CommentView struct {
	ID, WorkID, UserID, RootID, ParentID, ReplyToUserID int64
	Content, AuthorName, AuthorAvatar, ReplyToName      string
	ImageMediaID                                        int64
	ImageBucket, ImageObjectKey                         string
	CreatedAt                                           time.Time
	LikeCount, ReplyCount                               int64
}
