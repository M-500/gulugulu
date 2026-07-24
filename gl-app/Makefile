.PHONY: dev api media-worker

CONFIG_FILE ?= api/etc/gulu.yaml

# 本地开发时同时启动 API 和媒体消费者，避免视频任务只入队却无人处理。
dev:
	$(MAKE) -j2 api media-worker

api:
	go run ./api -f $(CONFIG_FILE)

media-worker:
	go run ./api/cmd/mediaworker -f $(CONFIG_FILE)
