package svc

import (
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"gl-app/api/internal/config"
	"gl-app/api/internal/models/user"
)

type ServiceContext struct {
	Config   config.Config
	UserRepo user.UserModel
}

func NewServiceContext(c config.Config) *ServiceContext {
	conn := sqlx.NewMysql(c.Mysql.DataSource)
	return &ServiceContext{
		Config:   c,
		UserRepo: user.NewUserModel(conn, c.CacheRedis),
	}
}
