package config

import (
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
		PresignExpire   int64
	}
}
