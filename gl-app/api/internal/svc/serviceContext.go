package svc

import (
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"gl-app/api/internal/config"
	"gl-app/api/internal/models/media"
	"gl-app/api/internal/models/user"
	workmodel "gl-app/api/internal/models/work"
	mediaqueue "gl-app/api/internal/queue"
)

type ServiceContext struct {
	Config          config.Config
	UserRepo        user.UserModel
	MediaAssetRepo  media.MediaAssetModel
	WorkRepo        workmodel.WorkModel
	WorkAssetRepo   workmodel.WorkAssetModel
	TopicRepo       workmodel.TopicModel
	WorkTopicRepo   workmodel.WorkTopicModel
	ProcessTaskRepo workmodel.MediaProcessTaskModel
	SqlConn         sqlx.SqlConn
	MinioClient     *minio.Client
	MediaQueue      *mediaqueue.MediaQueue
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
		Config:          c,
		UserRepo:        user.NewUserModel(conn, c.CacheRedis),
		MediaAssetRepo:  media.NewMediaAssetModel(conn, c.CacheRedis),
		WorkRepo:        workmodel.NewWorkModel(conn, c.CacheRedis),
		WorkAssetRepo:   workmodel.NewWorkAssetModel(conn, c.CacheRedis),
		TopicRepo:       workmodel.NewTopicModel(conn, c.CacheRedis),
		WorkTopicRepo:   workmodel.NewWorkTopicModel(conn, c.CacheRedis),
		ProcessTaskRepo: workmodel.NewMediaProcessTaskModel(conn, c.CacheRedis),
		SqlConn:         conn,
		MinioClient:     minioClient,
		MediaQueue:      mediaqueue.NewMediaQueue(c.MediaQueue.Brokers, c.MediaQueue.Topic),
	}
}
