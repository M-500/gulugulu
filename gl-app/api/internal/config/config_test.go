package config

import (
	"testing"

	"gl-app/pkg/gormx"
	"gl-app/pkg/redisx"

	"github.com/zeromicro/go-zero/core/conf"
)

func TestInfrastructureConfigCanBeLoadedFromYAML(t *testing.T) {
	content := []byte(`
Mysql:
  DataSource: user:pass@tcp(mysql:3306)/gulugulu
CacheRedis:
  - Host: redis:6379
`)
	var value struct {
		Mysql      gormx.MySQLConfig
		CacheRedis redisx.Config
	}
	if err := conf.LoadFromYamlBytes(content, &value); err != nil {
		t.Fatalf("LoadFromYamlBytes() error = %v", err)
	}
	if value.Mysql.MaxIdleConns != 10 || value.Mysql.MaxOpenConns != 100 || value.Mysql.ConnMaxLife != 3600 {
		t.Fatalf("MySQL默认连接池配置未生效: %+v", value.Mysql)
	}
	if len(value.CacheRedis) != 1 || value.CacheRedis[0].Type != "node" {
		t.Fatalf("Redis默认节点类型未生效: %+v", value.CacheRedis)
	}
}
