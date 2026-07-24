package config

import (
	"github.com/zeromicro/go-queue/kq"
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/rest"
)

type Config struct {
	rest.RestConf

	Mysql struct {
		DataSource string
	}

	CacheRedis cache.CacheConf

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
		PresignExpire   int64
	}

	MediaQueue struct {
		kq.KqConf
		MaxRetry int
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
