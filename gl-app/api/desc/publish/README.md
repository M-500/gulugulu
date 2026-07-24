# 作品发布接口

代码骨架由以下命令生成：

```bash
cd gl-app/api/desc
goctl api go -api bff.api -dir ../ --style=goZero

cd ../internal/models/media
goctl model mysql ddl -src media_asset.sql -dir . -c --style goZero

cd ../work
goctl model mysql ddl -src work.sql -dir . -c --style goZero
```

## 创建作品

`POST /api/v1/works` 使用 `multipart/form-data`：

- `payload`：JSON 字符串
- `cover`：jpg/jpeg/png 封面，最大 100MB
- `Idempotency-Key`：必填请求头，最多 64 字符

示例 `payload`：

```json
{
  "type": "video",
  "title": "测试视频",
  "content": "正文 #开发",
  "visibility": {
    "type": "public",
    "userIds": []
  },
  "topics": [
    {
      "id": 0,
      "name": "开发"
    }
  ],
  "assets": [
    {
      "mediaId": 10001,
      "sort": 0
    }
  ],
  "collectionId": 0,
  "original": true,
  "scheduledAt": null
}
```

## 状态流转

```text
pending -> processing -> succeeded
                           |
                           v
                    pending_review
                     /           \
                 rejected       approved
                                  |
                       published / scheduled
```

媒体处理成功不会直接发布。只有审核接口返回通过后，作品才会立即发布或进入定时发布。

创建作品接口只等待数据库事务和Kafka消息投递成功，不等待图片迁移、视频转码或审核。
前端收到 `workId` 后即可结束发布流程。

## 作品管理列表

```http
GET /api/v1/creator/works?page=1&pageSize=20
```

该接口只返回当前用户 `process_status=succeeded` 的作品，不返回媒体处理中的作品。
列表包含图片/视频类型、封面、审核状态、发布状态、审核失败原因和发布时间。

## 启动

先执行数据库迁移：

```text
api/internal/models/migrations/001_publish_work.sql
```

先启动 Redis（仅用于 go-zero 数据缓存）、Kafka 和 MinIO。媒体任务不再使用 Redis，
生产和消费均使用 go-zero `kq`：

```bash
cd deploy
docker compose up -d redis kafka minio minio-init
```

API 和 Worker 必须分别启动：

```bash
go run ./api -f api/etc/gulu.yaml
go run ./api/cmd/mediaworker -f api/etc/gulu.yaml
```

Worker 所在环境必须安装 `ffmpeg` 和 `ffprobe`。
