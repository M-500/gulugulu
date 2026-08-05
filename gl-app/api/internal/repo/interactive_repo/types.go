package interactive_repo

import (
	"gl-app/api/internal/constants"

	"gorm.io/gorm"
)

type InteractiveModel struct {
	gorm.Model
	ResourceID      int64             `gorm:"column:resource_id;not null;comment:资源ID;type:bigint;uniqueIndex:idx_resource_id_resource_type"`
	ResourceType    constants.BizType `gorm:"column:resource_type;not null;comment:资源类型;type:varchar(32);uniqueIndex:idx_resource_id_resource_type"`
	LikeCount       int64             `gorm:"column:like_count;default:0;comment:点赞数;type:bigint"`
	CommentCount    int64             `gorm:"column:comment_count;default:0;comment:评论数;type:bigint"`
	CollectionCount int64             `gorm:"column:collection_count;default:0;comment:收藏数;type:bigint"`
	ShareCount      int64             `gorm:"column:share_count;default:0;comment:分享数;type:bigint"`
	ViewCount       int64             `gorm:"column:view_count;default:0;comment:浏览数;type:bigint"`
}

func (InteractiveModel) TableName() string {
	return "interactive"
}
