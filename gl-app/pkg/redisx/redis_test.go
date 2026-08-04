package redisx

import (
	"testing"

	"github.com/redis/go-redis/v9"
)

func TestNew(t *testing.T) {
	t.Run("单节点", func(t *testing.T) {
		client, err := New(Config{{Host: "127.0.0.1:6379", Type: "node"}})
		if err != nil {
			t.Fatalf("New() error = %v", err)
		}
		defer client.Close()
		if _, ok := client.(*redis.Client); !ok {
			t.Fatalf("New() type = %T, want *redis.Client", client)
		}
	})

	t.Run("集群", func(t *testing.T) {
		client, err := New(Config{{Host: "redis-1:6379", Type: "cluster"}, {Host: "redis-2:6379", Type: "cluster"}})
		if err != nil {
			t.Fatalf("New() error = %v", err)
		}
		defer client.Close()
		if _, ok := client.(*redis.ClusterClient); !ok {
			t.Fatalf("New() type = %T, want *redis.ClusterClient", client)
		}
	})

	t.Run("空配置", func(t *testing.T) {
		if _, err := New(nil); err == nil {
			t.Fatal("New() error = nil, want error")
		}
	})
}
