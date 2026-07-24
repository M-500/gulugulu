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

视频作品的创建接口只等待数据库事务和Kafka消息投递成功，不等待视频转码或审核。
图片作品不进入Kafka，创建接口会通过MinIO服务端复制把正文图片和封面迁移到正式桶，
完成后直接进入 `pending_review`。前端收到 `workId` 后即可结束发布流程。

## 作品管理列表

```http
GET /api/v1/creator/works?page=1&pageSize=20&status=all&keyword=标题
```

该接口只返回当前用户 `process_status=succeeded` 的作品，不返回媒体处理中的作品。
列表包含图片/视频类型、封面、审核状态、发布状态、审核失败原因和发布时间。

`status` 支持：

- `all`：全部处理完成作品
- `published`：已经发布
- `reviewing`：等待审核
- `rejected`：审核未通过

作品卡片管理操作：

```http
PUT    /api/v1/creator/works/:workId/title
PUT    /api/v1/creator/works/:workId/visibility
DELETE /api/v1/creator/works/:workId
```

列表会额外返回视频时长、作品可见性，以及浏览、点赞、收藏和转发次数。

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
