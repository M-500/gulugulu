package comment_repo

import "gorm.io/gorm"

type CommentModel struct {
	gorm.Model
	ResourceID int64 `gorm:"comment:资源ID"`
}
