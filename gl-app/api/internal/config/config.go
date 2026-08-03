package config

import (
	"github.com/zeromicro/go-queue/kq"
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/rest"
)

type Config struct {
	rest.RestConf

	Mysql struct {
		DataSource   string
		MaxIdleConns int `json:",default=10"`
		MaxOpenConns int `json:",default=100"`
		ConnMaxLife  int `json:",default=3600"`
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
		PublicBucket    string
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
