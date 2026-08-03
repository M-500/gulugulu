package svc

import (
	"fmt"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	redisv9 "github.com/redis/go-redis/v9"
	"gorm.io/gorm"

	"gl-app/api/internal/config"
	mediaqueue "gl-app/api/internal/queue"
	mediarepo "gl-app/api/internal/repo/media_repo"
	userrepo "gl-app/api/internal/repo/user_repo"
	workrepo "gl-app/api/internal/repo/work_repo"
	"gl-app/pkg/gormx"
	"gl-app/pkg/redisx"
)

type ServiceContext struct {
	Config      config.Config
	GormDB      *gorm.DB
	RedisClient redisv9.UniversalClient
	UserRepo    userrepo.UserRepo
	WorkRepo    workrepo.WorkRepo
	MediaRepo   mediarepo.MediaRepo
	MinioClient *minio.Client
	MediaQueue  *mediaqueue.MediaQueue
}

func NewServiceContext(c config.Config) *ServiceContext {
	db, err := gormx.NewMySQL(c.Mysql)
	if err != nil {
		panic(err)
	}

	redisClient, err := redisx.New(c.CacheRedis)
	if err != nil {
		_ = gormx.Close(db)
		panic(err)
	}
	minioClient, err := minio.New(c.Minio.Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(c.Minio.AccessKeyID, c.Minio.SecretAccessKey, ""),
		Secure: c.Minio.UseSSL,
	})
	if err != nil {
		_ = redisx.Close(redisClient)
		_ = gormx.Close(db)
		panic(fmt.Errorf("初始化MinIO客户端失败: %w", err))
	}

	// 基础设施客户端只在 ServiceContext 创建一次，再通过构造函数注入各 Repo。
	userDAO := userrepo.NewUserDAO(db)
	userCache := userrepo.NewUserCache(redisClient)
	return &ServiceContext{
		Config:      c,
		GormDB:      db,
		RedisClient: redisClient,
		UserRepo:    userrepo.NewUserRepo(userDAO, userCache),
		WorkRepo:    workrepo.NewWorkRepo(db),
		MediaRepo:   mediarepo.NewMediaRepo(db),
		MinioClient: minioClient,
		MediaQueue:  mediaqueue.NewMediaQueue(c.MediaQueue.Brokers, c.MediaQueue.Topic),
	}
}

// Close 统一释放 ServiceContext 持有的长连接资源。
func (s *ServiceContext) Close() {
	if s.MediaQueue != nil {
		_ = s.MediaQueue.Close()
	}
	if s.RedisClient != nil {
		_ = redisx.Close(s.RedisClient)
	}
	if s.GormDB != nil {
		_ = gormx.Close(s.GormDB)
	}
}
