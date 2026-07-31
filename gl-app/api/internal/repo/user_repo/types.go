package user_repo

import (
	"time"

	"gorm.io/gorm"
)

type UserModel struct {
	gorm.Model
	NickName      string    `gorm:"type:varchar(64);not null;comment:'昵称'" json:"nick_name"`
	Email         string    `gorm:"type:varchar(64);not null;unique;comment:'邮箱'" json:"email"`
	Avatar        string    `gorm:"type:varchar(255);not null;comment:'头像'" json:"avatar"`
	Sex           int       `gorm:"type:int;not null;comment:'性别'" json:"sex"`
	LastLoginAt   time.Time `gorm:"type:datetime;not null;comment:'最后登录时间'" json:"last_login_at"`
	LastIpAddress string    `gorm:"type:varchar(64);not null;comment:'IP归属地'" json:"last_ip_address"`
	Age           int       `gorm:"type:int;not null;comment:'年龄'" json:"age"`
	GLID          string    `gorm:"type:varchar(64);not null;comment:'咕噜咕噜ID'" json:"glid"` // 这个是咕噜咕噜ID
	ProfileDesc   string    `gorm:"type:varchar(255);not null;comment:'个性描述'" json:"profile_desc"`
	Status        int       `gorm:"type:int;not null;comment:'状态'" json:"status"` // 0:正常 1:封禁 2:注销 3:异常

}

func (UserModel) TableName() string {
	return "gl_user"
}
