package user_repo

import (
	"time"

	"gorm.io/gorm"
)

// User 是 user 表对应的持久化实体。
// Repo 对外返回领域实体，业务层不再依赖 goctl 生成的 sqlx Model。
type User struct {
	ID          int64          `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	CreatedAt   time.Time      `gorm:"column:created_at" json:"createdAt"`
	UpdatedAt   time.Time      `gorm:"column:updated_at" json:"updatedAt"`
	DeletedAt   gorm.DeletedAt `gorm:"column:deleted_at;index" json:"-"`
	Email       string         `gorm:"column:email" json:"email"`
	Nickname    string         `gorm:"column:nickname" json:"nickname"`
	Avatar      string         `gorm:"column:avatar" json:"avatar"`
	Password    string         `gorm:"column:password" json:"-"`
	Sex         int64          `gorm:"column:sex" json:"sex"`
	LastLoginAt *time.Time     `gorm:"column:last_login_at" json:"lastLoginAt,omitempty"`
}

func (User) TableName() string {
	return "user"
}

// ErrNotFound 统一暴露仓储层的未找到错误，业务层无需感知 GORM 实现细节。
var ErrNotFound = gorm.ErrRecordNotFound
