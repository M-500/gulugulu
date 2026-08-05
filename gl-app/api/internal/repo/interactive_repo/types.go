package interactive_repo

import (
	"gl-app/api/internal/constants"
	"time"

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

// UserLikeModel 保存用户与资源之间的点赞关系，用版本号保证 Kafka 重投和乱序时幂等。
type UserLikeModel struct {
	ID           uint `gorm:"primaryKey"`
	CreatedAt    time.Time
	UpdatedAt    time.Time
	UserID       int64             `gorm:"column:user_id;not null;type:bigint;uniqueIndex:idx_user_resource_like"`
	ResourceID   int64             `gorm:"column:resource_id;not null;type:bigint;uniqueIndex:idx_user_resource_like"`
	ResourceType constants.BizType `gorm:"column:resource_type;not null;type:varchar(32);uniqueIndex:idx_user_resource_like"`
	Liked        bool              `gorm:"column:liked;not null;default:0;comment:当前是否点赞"`
	Version      int64             `gorm:"column:version;not null;default:0;comment:Redis操作版本"`
}

func (UserLikeModel) TableName() string { return "user_like" }

type LikeEvent struct {
	EventID      string            `json:"eventId"`
	UserID       int64             `json:"userId"`
	ResourceID   int64             `json:"resourceId"`
	ResourceType constants.BizType `json:"resourceType"`
	Liked        bool              `json:"liked"`
	Version      int64             `json:"version"`
	OccurredAt   time.Time         `json:"occurredAt"`
}

type LikeMutation struct {
	Changed bool
	Liked   bool
	Count   int64
	Version int64
}
