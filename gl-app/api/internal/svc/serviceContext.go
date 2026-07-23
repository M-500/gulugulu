package svc

import (
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"gl-app/api/internal/config"
	"gl-app/api/internal/models/media"
	"gl-app/api/internal/models/user"
)

type ServiceContext struct {
	Config         config.Config
	UserRepo       user.UserModel
	MediaAssetRepo media.MediaAssetModel
	MinioClient    *minio.Client
}

func NewServiceContext(c config.Config) *ServiceContext {
	conn := sqlx.NewMysql(c.Mysql.DataSource)
	minioClient, err := minio.New(c.Minio.Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(c.Minio.AccessKeyID, c.Minio.SecretAccessKey, ""),
		Secure: c.Minio.UseSSL,
	})
	if err != nil {
		panic(err)
	}

	return &ServiceContext{
		Config:         c,
		UserRepo:       user.NewUserModel(conn, c.CacheRedis),
		MediaAssetRepo: media.NewMediaAssetModel(conn, c.CacheRedis),
		MinioClient:    minioClient,
	}
}
