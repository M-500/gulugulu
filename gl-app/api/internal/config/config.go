package config

import (
	"gl-app/pkg/gormx"
	"gl-app/pkg/redisx"

	"github.com/zeromicro/go-queue/kq"
	"github.com/zeromicro/go-zero/rest"
)

type Config struct {
	rest.RestConf

	Mysql      gormx.MySQLConfig
	CacheRedis redisx.Config

	Salt string // 加密用的盐

	JwtAuth struct {
		AccessSecret string
		AccessExpire int64
	}

	Minio struct {
		Endpoint        string
		AccessKeyID     string
		SecretAccessKey string
		UseSSL          bool
		TempBucket      string
		FormalBucket    string
		PublicBucket    string
		PresignExpire   int64
	}

	MediaQueue struct {
		kq.KqConf
		MaxRetry int
	}

	LikeQueue struct {
		kq.KqConf
	}

	MediaWorker struct {
		FFmpegPath  string
		FFprobePath string
		TempDir     string
		HlsTime     int
	}

	Audit struct {
		AdminUserIds []int64
	}
}
