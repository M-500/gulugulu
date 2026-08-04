// Package redisx 统一创建 Redis 单机或集群客户端。
package redisx

import (
	"fmt"
	"strings"

	"github.com/redis/go-redis/v9"
)

type NodeConfig struct {
	Host string
	Type string `json:",default=node,options=node|cluster"`
	Pass string `json:",optional"`
}

type Config []NodeConfig

// New 根据节点数量和类型创建单机或集群客户端，不在初始化阶段强制访问网络。
// 这样 Redis 短暂不可用时，Repo 仍可按既定策略降级读取数据库。
func New(config Config) (redis.UniversalClient, error) {
	if len(config) == 0 {
		return nil, fmt.Errorf("Redis至少需要配置一个节点")
	}
	addresses := make([]string, 0, len(config))
	for _, node := range config {
		host := strings.TrimSpace(node.Host)
		if host == "" {
			return nil, fmt.Errorf("Redis节点地址不能为空")
		}
		addresses = append(addresses, host)
	}
	if len(addresses) > 1 || strings.EqualFold(config[0].Type, "cluster") {
		return redis.NewClusterClient(&redis.ClusterOptions{
			Addrs: addresses, Password: config[0].Pass,
		}), nil
	}
	return redis.NewClient(&redis.Options{
		Addr: addresses[0], Password: config[0].Pass,
	}), nil
}

func Close(client redis.UniversalClient) error {
	if client == nil {
		return nil
	}
	return client.Close()
}
