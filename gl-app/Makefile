.PHONY: dev api

CONFIG_FILE ?= api/etc/gulu.yaml

# API 已通过 servicegroup 集成媒体消费者，启动一个进程即可。
dev: api

api:
	go run ./api -f $(CONFIG_FILE)
