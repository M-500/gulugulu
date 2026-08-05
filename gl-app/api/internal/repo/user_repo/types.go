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
	Email       string         `gorm:"column:email;comment:用户邮箱;type:varchar(256);unique" json:"email"`
	Nickname    string         `gorm:"column:nickname;comment:用户昵称;type:varchar(256);not null" json:"nickname"`
	Avatar      string         `gorm:"column:avatar;comment:用户头像;type:varchar(256)" json:"avatar"`
	Password    string         `gorm:"column:password;comment:用户密码;type:varchar(256);not null	" json:"-"`
	Sex         int64          `gorm:"column:sex;comment:用户性别" json:"sex"`
	LastLoginAt *time.Time     `gorm:"column:last_login_at;comment:最后登录时间" json:"lastLoginAt,omitempty"`
	IPAddress   string         `gorm:"column:ip_address;comment:用户IP地址;type:varchar(64)" json:"ipAddress"`
	Bio         string         `gorm:"column:bio;comment:用户简介;type:varchar(256)" json:"bio"`
	IsAudit     bool           `gorm:"column:is_audit;comment:是否是审核员;not null;default:0" json:"isAudit"`
	BothDay     *time.Time     `gorm:"column:both_day;comment:用户生日" json:"bothDay,omitempty"`
}

func (User) TableName() string {
	return "user"
}

// ErrNotFound 统一暴露仓储层的未找到错误，业务层无需感知 GORM 实现细节。
var ErrNotFound = gorm.ErrRecordNotFound
