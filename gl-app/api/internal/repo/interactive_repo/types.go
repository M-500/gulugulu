package interactive_repo

import "gorm.io/gorm"

type InteractiveModel struct {
	gorm.Model
	ResourceID      int64 `gorm:"not null;comment:资源ID;type:bigint;index"`
	LikeCount       int64 `gorm:"default:0;comment:点赞数;type:bigint"`
	CommentCount    int64 `gorm:"default:0;comment:评论数;type:bigint"`
	CollectionCount int64 `gorm:"default:0;comment:收藏数;type:bigint"`
	ShareCount      int64 `gorm:"default:0;comment:分享数;type:bigint"`
	ViewCount       int64 `gorm:"default:0;comment:浏览数;type:bigint"`
}

func (InteractiveModel) TableName() string {
	return "interactive"
}
