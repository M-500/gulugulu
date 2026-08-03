package follow_repo

import "gorm.io/gorm"

// 关注相关的

type FollowModel struct {
	gorm.Model
	FollowerUserID int64 `gorm:"comment:被关注者ID;type:bigint;not null;index"`
	FansID         int64 `gorm:"comment:关注人的ID;type:bigint;not null;index"`
}
