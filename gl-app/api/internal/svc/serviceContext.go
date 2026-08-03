package svc

import (
	"fmt"
	"time"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	redisv9 "github.com/redis/go-redis/v9"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"gl-app/api/internal/config"
	mediaqueue "gl-app/api/internal/queue"
	mediarepo "gl-app/api/internal/repo/media_repo"
	userrepo "gl-app/api/internal/repo/user_repo"
	workrepo "gl-app/api/internal/repo/work_repo"
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
	db, err := gorm.Open(mysql.Open(c.Mysql.DataSource), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Warn),
	})
	if err != nil {
		panic(fmt.Errorf("初始化GORM失败: %w", err))
	}
	configureConnectionPool(db, c)

	redisClient := newRedisClient(c)
	minioClient, err := minio.New(c.Minio.Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(c.Minio.AccessKeyID, c.Minio.SecretAccessKey, ""),
		Secure: c.Minio.UseSSL,
	})
	if err != nil {
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

func configureConnectionPool(db *gorm.DB, c config.Config) {
	sqlDB, err := db.DB()
	if err != nil {
		panic(fmt.Errorf("获取GORM底层连接池失败: %w", err))
	}
	sqlDB.SetMaxIdleConns(c.Mysql.MaxIdleConns)
	sqlDB.SetMaxOpenConns(c.Mysql.MaxOpenConns)
	sqlDB.SetConnMaxLifetime(time.Duration(c.Mysql.ConnMaxLife) * time.Second)
}

func newRedisClient(c config.Config) redisv9.UniversalClient {
	if len(c.CacheRedis) == 0 {
		panic("CacheRedis至少需要配置一个Redis节点")
	}
	addresses := make([]string, 0, len(c.CacheRedis))
	for _, node := range c.CacheRedis {
		addresses = append(addresses, node.Host)
	}
	if len(addresses) > 1 || c.CacheRedis[0].Type == "cluster" {
		return redisv9.NewClusterClient(&redisv9.ClusterOptions{
			Addrs:    addresses,
			Password: c.CacheRedis[0].Pass,
		})
	}
	return redisv9.NewClient(&redisv9.Options{
		Addr:     addresses[0],
		Password: c.CacheRedis[0].Pass,
	})
}

// Close 统一释放 ServiceContext 持有的长连接资源。
func (s *ServiceContext) Close() {
	if s.MediaQueue != nil {
		_ = s.MediaQueue.Close()
	}
	if s.RedisClient != nil {
		_ = s.RedisClient.Close()
	}
	if s.GormDB != nil {
		if sqlDB, err := s.GormDB.DB(); err == nil {
			_ = sqlDB.Close()
		}
	}
}
