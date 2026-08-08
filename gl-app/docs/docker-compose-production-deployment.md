# Gulugulu 单机 Docker Compose 生产部署指南

> 适用范围：不使用 Jenkins，在一台 Linux 服务器上通过 Docker Compose 部署 Gulugulu App 前台、CMS 管理端、Go API、MySQL、Redis、Kafka、MinIO，并由 Caddy 提供统一公网入口和自动 HTTPS。
>
> 本文按当前仓库结构编写。命令默认在仓库根目录或 `gl-app/deploy` 目录执行，请勿直接照搬仓库中的开发密码、IP 地址和 `localhost` 配置到生产环境。

## 1. 部署目标与原则

推荐的单机生产拓扑如下：

```text
Internet
   |
   | 80 / 443
   v
Caddy（TLS、域名入口）
   |-- app.example.com   -> gulugulu-app:80
   |-- cms.example.com   -> gulugulu-cms:80
   `-- media.example.com -> minio:9000

gulugulu-app / gulugulu-cms
   |-- /app、/api、/na -> gulugulu-api:8888
   `-- 静态资源由各自 Nginx 提供

gulugulu-api
   |-- mysql:3306
   |-- redis:6379
   |-- kafka:9092
   `-- media.example.com:443 -> Caddy -> minio:9000
```

生产部署遵循以下原则：

1. 只有 Caddy 的 `80/443` 端口暴露到公网。
2. API、MySQL、Redis、Kafka、MinIO API 仅加入 `gulugulu-network`，不直接暴露公网端口。
3. App、CMS 与 API 使用同域反向代理，避免额外的 API CORS 配置。
4. MinIO 使用独立 HTTPS 域名。后端签发的上传和下载 URL 必须是浏览器可访问的公网地址，不能是 `minio:9000`、`localhost:9000` 等容器内或本机地址。
5. 生产密码、JWT Secret、加密盐不写入 Git，不复用仓库中的开发值。
6. 每次发布使用 Git Commit 等不可变镜像标签，不依赖 `latest`，以便快速回滚。
7. 数据库和对象存储必须同时备份；备份必须定期做恢复演练。

## 2. 当前仓库中的部署组件

| 组件 | 构建文件 | 运行端口 | 说明 |
| --- | --- | ---: | --- |
| App 前台 | `front/app/Dockerfile` | 80 | Vite 构建，Nginx 提供静态文件和 API 代理 |
| CMS 管理端 | `front/cms/Dockerfile` | 80 | Vue CLI 构建，Nginx 提供静态文件和 API 代理 |
| Go API | `gl-app/Dockerfile` | 8888 | 同一进程运行 HTTP API、媒体 Worker 和点赞 Worker |
| MySQL | Compose 镜像 | 3306 | 业务主数据 |
| Redis | Compose 镜像 | 6379 | 缓存 |
| Kafka | Compose 镜像 | 9092 | 媒体处理与点赞事件队列 |
| MinIO | Compose 镜像 | 9000/9001 | 对象存储 API/管理控制台 |
| Caddy | Compose 镜像 | 80/443 | 公网入口、自动签发和续期 TLS 证书 |

仓库已有的 `gl-app/deploy/docker-compose.yaml` 更接近本地开发基础设施配置，存在以下生产问题：

- 包含硬编码密码和本地地址。
- MySQL、Redis、Kafka、MinIO 端口对宿主机公开。
- 没有运行 App、CMS 和 API 容器。
- `MINIO_SERVER_URL` 使用 `localhost`，远程浏览器无法访问。
- `gl-app/api/etc/gulu.yaml` 使用局域网 IP、`localhost` 和开发 Secret。
- `gl-app/deploy/conf/my.cnf` 的 `innodb_buffer_pool_size` 是 `24G`，但现有 Compose 给 MySQL 的内存上限是 `8G`，两者冲突，可能导致 OOM。除非按服务器内存重新调优，否则生产环境不要直接挂载这个文件。

因此推荐保留开发 Compose，另建生产专用的 `docker-compose.prod.yaml`、`.env.production`、`Caddyfile` 和 `prod/gulu.yaml`。

## 3. 服务器与域名前置条件

### 3.1 推荐服务器规格

最低可运行规格取决于视频并发和素材大小。建议起步配置：

| 资源 | 建议 |
| --- | --- |
| CPU | 8 核及以上；FFmpeg 转码会持续占用 CPU |
| 内存 | 16 GB 及以上；中等负载建议 32 GB |
| 系统盘 | 80 GB SSD 及以上 |
| 数据盘 | 按媒体量规划，建议独立 SSD/云盘并监控使用率 |
| 操作系统 | Ubuntu 24.04 LTS 或 Debian 12 |
| 架构 | 当前后端 Dockerfile 固定构建 `GOARCH=amd64`，服务器应使用 x86_64/amd64 |

如果服务器是 ARM64，必须先让 `gl-app/Dockerfile` 的 `GOARCH` 可通过构建参数配置，并构建 ARM64 镜像；否则可能出现 `exec format error`。

### 3.2 域名和 DNS

准备三个域名，并都添加指向服务器公网 IP 的 A/AAAA 记录：

```text
app.example.com    -> 服务器公网 IP
cms.example.com    -> 服务器公网 IP
media.example.com  -> 服务器公网 IP
```

DNS 生效后再启动 Caddy。Let's Encrypt/ZeroSSL 必须能够从公网访问服务器的 80 和 443 端口。

### 3.3 防火墙与安全组

公网只开放：

| 端口 | 用途 | 来源限制 |
| ---: | --- | --- |
| 22 | SSH | 强烈建议只允许管理 IP，使用密钥登录 |
| 80 | HTTP/ACME | 公网；Caddy 会跳转 HTTPS |
| 443 | HTTPS | 公网 |

不要对公网开放 `3306`、`6379`、`8888`、`9000`、`9092`、`9093`、`9001`。MinIO 控制台通过 SSH 隧道访问，见后文。

### 3.4 安装 Docker

使用 Docker 官方仓库安装 Docker Engine 和 Compose Plugin。安装后确认：

```bash
docker version
docker compose version
```

建议 Docker Engine 26+、Compose v2.24+。部署账号应有执行 Docker 的权限。注意：加入 `docker` 组等同于拥有宿主机 root 权限，应严格限制账号和 SSH 密钥。

## 4. 生产目录规划

本文采用“服务器上保留源码并本地构建镜像”的方式：

```text
/opt/gulugulu/
├── source/                         # Git 仓库
│   ├── front/
│   └── gl-app/
│       └── deploy/
│           ├── docker-compose.prod.yaml
│           ├── .env.production     # 不提交 Git，权限 600
│           ├── Caddyfile
│           └── prod/
│               └── gulu.yaml       # 不提交 Git，权限 600
└── backup/                         # 备份目录，建议位于独立数据盘
    ├── mysql/
    └── minio/
```

克隆代码：

```bash
sudo install -d -o "$USER" -g "$USER" /opt/gulugulu/source /opt/gulugulu/backup
git clone <你的仓库地址> /opt/gulugulu/source
cd /opt/gulugulu/source
```

生产机建议只允许从受保护的 release tag 或明确 commit 部署，不要直接发布未经确认的工作区代码。

## 5. 准备生产环境变量

在 `gl-app/deploy/.env.production` 中保存 Compose 使用的非代码配置：

```dotenv
# Compose 项目与镜像标签。每次发布将 IMAGE_TAG 改成 Git 短提交号或版本号。
COMPOSE_PROJECT_NAME=gulugulu
IMAGE_TAG=2026.08.07-abcdef123456

# 域名，不带协议和路径。
APP_DOMAIN=app.example.com
CMS_DOMAIN=cms.example.com
MEDIA_DOMAIN=media.example.com
ACME_EMAIL=ops@example.com

# 数据库。生产环境请使用随机值，禁止照抄示例。
MYSQL_DATABASE=gulugulu
MYSQL_USER=gulugulu
MYSQL_PASSWORD=REPLACE_WITH_RANDOM_DATABASE_PASSWORD
MYSQL_ROOT_PASSWORD=REPLACE_WITH_RANDOM_ROOT_PASSWORD

# Redis。
REDIS_PASSWORD=REPLACE_WITH_RANDOM_REDIS_PASSWORD

# MinIO。用户名不要使用默认 admin，密码至少 20 位。
MINIO_ROOT_USER=gulugulu-storage-admin
MINIO_ROOT_PASSWORD=REPLACE_WITH_RANDOM_MINIO_PASSWORD

# 时区。
TZ=Asia/Shanghai
```

设置文件权限：

```bash
cd /opt/gulugulu/source/gl-app/deploy
chmod 600 .env.production
```

生成随机值时可以使用：

```bash
openssl rand -hex 32
```

注意事项：

- `.env.production` 必须加入 `.gitignore`，不得提交。
- Compose 环境变量最终可被有 Docker 权限的用户查看。更高安全等级的环境应接入 Vault、云 Secret Manager 或 Docker Secrets，而不是长期使用 `.env`。
- 不要在密码中使用换行符。为降低 DSN 转义风险，数据库密码可直接使用 `openssl rand -hex 32` 生成的十六进制字符串。
- 不要执行会把完整 Compose 渲染结果上传到日志平台的命令，因为渲染结果可能包含 Secret。

## 6. 准备后端生产配置

创建 `gl-app/deploy/prod/gulu.yaml`：

```yaml
Name: gulu
Host: 0.0.0.0
Port: 8888
Mode: pro
MaxBytes: 104857600
Timeout: 300000

Mysql:
  # 密码必须与 .env.production 一致；Compose 不会自动替换挂载 YAML 中的变量。
  DataSource: gulugulu:REPLACE_WITH_RANDOM_DATABASE_PASSWORD@tcp(mysql:3306)/gulugulu?charset=utf8mb4&parseTime=true&loc=Asia%2FShanghai
  MaxIdleConns: 10
  MaxOpenConns: 50
  ConnMaxLife: 3600

CacheRedis:
  - Host: redis:6379
    Type: node
    Pass: REPLACE_WITH_RANDOM_REDIS_PASSWORD

# 用 openssl rand -hex 32 分别生成，禁止复用仓库开发配置中的值。
Salt: REPLACE_WITH_RANDOM_PASSWORD_SALT

JwtAuth:
  AccessSecret: REPLACE_WITH_RANDOM_JWT_SECRET
  # 示例为 7 天。生产环境不建议使用仓库中 1 年的超长有效期。
  AccessExpire: 604800

Minio:
  # 必须是浏览器可访问的 HTTPS 域名，不能写 minio:9000 或 localhost:9000。
  # Caddy 在 Docker 网络内注册同名 alias，API 会通过 HTTPS 访问 Caddy。
  Endpoint: media.example.com
  AccessKeyID: gulugulu-storage-admin
  SecretAccessKey: REPLACE_WITH_RANDOM_MINIO_PASSWORD
  UseSSL: true
  TempBucket: gulugulu-temp
  FormalBucket: gulugulu-media
  PublicBucket: gulugulu-public
  PresignExpire: 900

MediaQueue:
  Name: media-worker
  Brokers:
    - kafka:9092
  Group: gulugulu-media-workers
  Topic: gulugulu-media-tasks
  Offset: first
  Conns: 1
  Consumers: 1
  Processors: 2
  ForceCommit: false
  CommitInOrder: false
  MaxRetry: 3

MediaWorker:
  FFmpegPath: ffmpeg
  FFprobePath: ffprobe
  TempDir: /tmp/gulugulu-media
  HlsTime: 6

LikeQueue:
  Name: like-worker
  Brokers:
    - kafka:9092
  Group: gulugulu-like-workers
  Topic: gulugulu-like-events
  Offset: first
  Conns: 1
  Consumers: 1
  Processors: 4
  ForceCommit: false
  CommitInOrder: true

# 当前代码仍使用用户 ID 白名单。上线前应确认管理员 ID；长期应改成 RBAC。
Audit:
  AdminUserIds: [1]
```

必须将以下值替换为真实值：

- 数据库密码；
- Redis 密码；
- MinIO 用户名和密码；
- `Salt`；
- `JwtAuth.AccessSecret`；
- `Minio.Endpoint` 域名；
- 审核管理员 ID。

设置权限：

```bash
mkdir -p prod
chmod 700 prod
chmod 600 prod/gulu.yaml
```

### 6.1 为什么 MinIO Endpoint 要使用公网域名

后端会基于 `Minio.Endpoint` 生成预签名上传、下载 URL 和公共头像 URL。如果配置为 `minio:9000`，这个名字只有 Docker 网络中的容器能解析，用户浏览器无法访问；如果配置为 `localhost:9000`，浏览器访问的是用户自己的电脑。

本文让 `media.example.com` 在公网解析到服务器，在 Docker 网络内部则作为 Caddy 的网络别名。这样：

- 浏览器访问 `https://media.example.com`；
- API 容器也访问 `https://media.example.com`；
- Caddy 将流量转发到 `minio:9000`；
- 预签名 URL 的 Host 和协议在内外网保持一致，签名不会因为改写 Host 而失效。

## 7. 准备 Caddy 配置

创建 `gl-app/deploy/Caddyfile`：

```caddyfile
{
    email {$ACME_EMAIL}
    admin off
}

(security_headers) {
    header {
        Strict-Transport-Security "max-age=31536000; includeSubDomains"
        X-Content-Type-Options "nosniff"
        Referrer-Policy "strict-origin-when-cross-origin"
        -Server
    }
}

{$APP_DOMAIN} {
    import security_headers
    encode zstd gzip
    reverse_proxy app:80
}

{$CMS_DOMAIN} {
    import security_headers
    encode zstd gzip
    reverse_proxy cms:80
}

{$MEDIA_DOMAIN} {
    import security_headers

    # 与 API、前端及 Nginx 的 100 MB 限制保持一致。
    request_body {
        max_size 100MB
    }

    reverse_proxy minio:9000
}
```

Caddy 会自动申请、保存并续期证书。`caddy_data` 卷不能随意删除，否则可能触发 CA 频率限制。

生产环境修改 Caddyfile 后，优先无损重载：

```bash
docker compose --env-file .env.production -f docker-compose.prod.yaml exec caddy \
  caddy reload --config /etc/caddy/Caddyfile
```

## 8. 推荐的生产 Docker Compose

创建 `gl-app/deploy/docker-compose.prod.yaml`：

```yaml
name: ${COMPOSE_PROJECT_NAME:-gulugulu}

services:
  mysql:
    image: mysql:8.0
    restart: unless-stopped
    stop_grace_period: 60s
    environment:
      TZ: ${TZ:-Asia/Shanghai}
      MYSQL_DATABASE: ${MYSQL_DATABASE}
      MYSQL_USER: ${MYSQL_USER}
      MYSQL_PASSWORD: ${MYSQL_PASSWORD}
      MYSQL_ROOT_PASSWORD: ${MYSQL_ROOT_PASSWORD}
    command:
      - --character-set-server=utf8mb4
      - --collation-server=utf8mb4_0900_ai_ci
      - --default-time-zone=+08:00
      - --max-connections=300
      - --innodb-buffer-pool-size=2G
      - --innodb-flush-log-at-trx-commit=1
      - --sync-binlog=1
    volumes:
      - mysql_data:/var/lib/mysql
    healthcheck:
      test: ["CMD-SHELL", "mysqladmin ping -h 127.0.0.1 -uroot -p\"$$MYSQL_ROOT_PASSWORD\" --silent"]
      interval: 10s
      timeout: 5s
      retries: 20
      start_period: 40s
    networks:
      - gulugulu-network

  redis:
    image: bitnami/redis:7.2
    restart: unless-stopped
    environment:
      REDIS_PASSWORD: ${REDIS_PASSWORD}
      REDIS_AOF_ENABLED: "yes"
      TZ: ${TZ:-Asia/Shanghai}
    volumes:
      - redis_data:/bitnami/redis/data
    healthcheck:
      test: ["CMD-SHELL", "redis-cli -a \"$$REDIS_PASSWORD\" ping | grep PONG"]
      interval: 10s
      timeout: 5s
      retries: 10
      start_period: 20s
    networks:
      - gulugulu-network

  kafka:
    image: bitnami/kafka:3.9.0
    restart: unless-stopped
    stop_grace_period: 60s
    ulimits:
      nofile:
        soft: 65536
        hard: 65536
    environment:
      TZ: ${TZ:-Asia/Shanghai}
      KAFKA_CFG_NODE_ID: 0
      KAFKA_CFG_PROCESS_ROLES: controller,broker
      KAFKA_CFG_CONTROLLER_QUORUM_VOTERS: 0@kafka:9093
      KAFKA_CFG_LISTENERS: PLAINTEXT://:9092,CONTROLLER://:9093
      KAFKA_CFG_ADVERTISED_LISTENERS: PLAINTEXT://kafka:9092
      KAFKA_CFG_LISTENER_SECURITY_PROTOCOL_MAP: CONTROLLER:PLAINTEXT,PLAINTEXT:PLAINTEXT
      KAFKA_CFG_CONTROLLER_LISTENER_NAMES: CONTROLLER
      KAFKA_CFG_AUTO_CREATE_TOPICS_ENABLE: "true"
      KAFKA_CFG_NUM_PARTITIONS: 3
      KAFKA_CFG_LOG_RETENTION_HOURS: 168
    volumes:
      - kafka_data:/bitnami/kafka
    healthcheck:
      test: ["CMD-SHELL", "/opt/bitnami/kafka/bin/kafka-topics.sh --bootstrap-server 127.0.0.1:9092 --list >/dev/null 2>&1"]
      interval: 15s
      timeout: 10s
      retries: 20
      start_period: 40s
    networks:
      - gulugulu-network

  minio:
    image: minio/minio:RELEASE.2025-01-20T14-49-07Z
    restart: unless-stopped
    command: server /data --console-address ":9001"
    environment:
      TZ: ${TZ:-Asia/Shanghai}
      MINIO_ROOT_USER: ${MINIO_ROOT_USER}
      MINIO_ROOT_PASSWORD: ${MINIO_ROOT_PASSWORD}
      MINIO_SERVER_URL: https://${MEDIA_DOMAIN}
      MINIO_BROWSER: "on"
      MINIO_API_CORS_ALLOW_ORIGIN: https://${APP_DOMAIN},https://${CMS_DOMAIN}
    volumes:
      - minio_data:/data
    # 控制台只监听宿主机回环地址，通过 SSH 隧道访问，绝不直接开放公网。
    ports:
      - "127.0.0.1:19001:9001"
    healthcheck:
      test: ["CMD", "curl", "-f", "http://127.0.0.1:9000/minio/health/live"]
      interval: 15s
      timeout: 5s
      retries: 10
      start_period: 20s
    networks:
      - gulugulu-network

  api:
    image: gulugulu-api:${IMAGE_TAG}
    build:
      context: ..
      dockerfile: Dockerfile
      args:
        VERSION: ${IMAGE_TAG}
        GIT_COMMIT: ${IMAGE_TAG}
    restart: unless-stopped
    init: true
    stop_grace_period: 90s
    environment:
      TZ: ${TZ:-Asia/Shanghai}
    volumes:
      - ./prod/gulu.yaml:/app/etc/gulu.yaml:ro
      - media_tmp:/tmp/gulugulu-media
    depends_on:
      mysql:
        condition: service_healthy
      redis:
        condition: service_healthy
      kafka:
        condition: service_healthy
      minio:
        condition: service_healthy
      # API 启动时 Media Worker 会立即通过媒体域名检查/创建 Bucket；
      # 必须等 Caddy 已取得证书且媒体路由可用，否则 API 会启动失败。
      caddy:
        condition: service_healthy
    healthcheck:
      test: ["CMD", "curl", "-fsS", "http://127.0.0.1:8888/app/v1/recommend/works?page=1&pageSize=1"]
      interval: 30s
      timeout: 5s
      retries: 5
      start_period: 60s
    networks:
      - gulugulu-network

  app:
    image: gulugulu-app:${IMAGE_TAG}
    build:
      context: ../../front/app
      dockerfile: Dockerfile
      args:
        # 留空以使用同域 Nginx 反向代理，避免浏览器跨域。
        VITE_API_BASE_URL: ""
    restart: unless-stopped
    environment:
      API_UPSTREAM: http://api:8888
      MINIO_UPSTREAM: http://minio:9000
    depends_on:
      api:
        condition: service_healthy
    healthcheck:
      test: ["CMD", "wget", "-q", "-O", "/dev/null", "http://127.0.0.1/healthz"]
      interval: 30s
      timeout: 5s
      retries: 3
      start_period: 10s
    networks:
      - gulugulu-network

  cms:
    image: gulugulu-cms:${IMAGE_TAG}
    build:
      context: ../../front/cms
      dockerfile: Dockerfile
      args:
        VUE_APP_API_BASE_URL: ""
    restart: unless-stopped
    environment:
      API_UPSTREAM: http://api:8888
      MINIO_UPSTREAM: http://minio:9000
    depends_on:
      api:
        condition: service_healthy
    healthcheck:
      test: ["CMD", "wget", "-q", "-O", "/dev/null", "http://127.0.0.1/healthz"]
      interval: 30s
      timeout: 5s
      retries: 3
      start_period: 10s
    networks:
      - gulugulu-network

  caddy:
    image: caddy:2.8-alpine
    restart: unless-stopped
    environment:
      APP_DOMAIN: ${APP_DOMAIN}
      CMS_DOMAIN: ${CMS_DOMAIN}
      MEDIA_DOMAIN: ${MEDIA_DOMAIN}
      ACME_EMAIL: ${ACME_EMAIL}
    ports:
      - "80:80"
      - "443:443"
      - "443:443/udp"
    volumes:
      - ./Caddyfile:/etc/caddy/Caddyfile:ro
      - caddy_data:/data
      - caddy_config:/config
    depends_on:
      minio:
        condition: service_healthy
    healthcheck:
      # 不只校验配置，还要等待证书签发成功且 MinIO HTTPS 路由真正可用。
      test: ["CMD-SHELL", "wget -q --spider https://${MEDIA_DOMAIN}/minio/health/live"]
      interval: 10s
      timeout: 10s
      retries: 30
      start_period: 20s
    networks:
      gulugulu-network:
        aliases:
          # 让 API 容器将媒体公网域名解析到 Caddy，保持预签名 URL 的 Host/协议一致。
          - ${MEDIA_DOMAIN}

volumes:
  mysql_data:
  redis_data:
  kafka_data:
  minio_data:
  media_tmp:
  caddy_data:
  caddy_config:

networks:
  gulugulu-network:
    name: gulugulu-network
    driver: bridge
```

### 8.1 关于镜像版本固定

示例固定了 MinIO 和 Kafka 的明确版本，但 MySQL、Redis、Caddy 仍是 minor tag。更严格的生产环境应将所有镜像固定到经过验证的 patch 版本，甚至固定镜像 digest：

```yaml
image: caddy:2.8.x-alpine@sha256:<verified-digest>
```

升级基础镜像前先在预发布环境验证数据兼容性，不要使用无版本的 `latest` 自动跨大版本升级。

### 8.2 关于资源限制

不要盲目复制固定资源限制。应先观察一段时间的 CPU、内存、磁盘 IO 和视频任务耗时，再按服务器容量设置。特别注意：

- FFmpeg 的峰值 CPU 和临时磁盘占用可能很高；
- MySQL Buffer Pool 通常可设为分配给 MySQL 内存的 50%～70%，但同机部署必须给其他容器和系统留足内存；
- Kafka 与 MinIO 都依赖稳定磁盘，生产数据盘不应接近满盘；
- Docker Compose 的 `deploy.resources` 在不同运行方式下语义容易被误解，单机 Compose 可优先使用受当前 Compose 版本支持的 `mem_limit`、`cpus`，设置后必须实测。

## 9. 首次部署

### 9.1 发布前检查

进入部署目录：

```bash
cd /opt/gulugulu/source/gl-app/deploy
```

检查以下事项：

```bash
# 文件存在且权限正确
ls -l .env.production prod/gulu.yaml Caddyfile docker-compose.prod.yaml

# 确认域名已解析到当前服务器
getent hosts app.example.com
getent hosts cms.example.com
getent hosts media.example.com

# 确认 80/443 未被其他程序占用
sudo ss -lntup | grep -E ':(80|443)\b' || true

# 查看磁盘和内存
df -h
free -h
```

校验 Compose。该命令的输出包含渲染后的配置，可能含 Secret，因此只做静默校验：

```bash
docker compose --env-file .env.production -f docker-compose.prod.yaml config --quiet
```

校验 Caddyfile：

```bash
docker run --rm \
  --env-file .env.production \
  -v "$PWD/Caddyfile:/etc/caddy/Caddyfile:ro" \
  caddy:2.8-alpine caddy validate --config /etc/caddy/Caddyfile
```

### 9.2 构建镜像

先把 `IMAGE_TAG` 设置为当前 Git Commit：

```bash
git -C /opt/gulugulu/source rev-parse --short=12 HEAD
```

将输出写入 `.env.production` 的 `IMAGE_TAG`，然后构建：

```bash
docker compose --env-file .env.production -f docker-compose.prod.yaml build --pull api app cms
```

构建前建议执行测试：

```bash
# Go 后端
cd /opt/gulugulu/source/gl-app
go test ./...

# App
cd /opt/gulugulu/source/front/app
npm ci
npm run lint
npm run build

# CMS
cd /opt/gulugulu/source/front/cms
npm ci
npm run lint
npm run build
```

如果不希望在生产机安装 Go/Node，可用 Dockerfile 的构建阶段做检查，或在一台独立构建机生成带版本号的镜像并推送到私有仓库；生产 Compose 只执行 `docker compose pull && docker compose up -d`。这仍然不需要 Jenkins，并且比在生产机编译更符合职责分离原则。

### 9.3 分阶段启动

先启动基础设施和 Caddy。Caddy 此时访问 App/CMS 会短暂返回上游不可用，但会先完成媒体域名证书签发，为 API 启动时的 MinIO 检查做好准备：

```bash
cd /opt/gulugulu/source/gl-app/deploy
docker compose --env-file .env.production -f docker-compose.prod.yaml up -d mysql redis kafka minio caddy
docker compose --env-file .env.production -f docker-compose.prod.yaml ps
```

等待全部 healthy，再启动 API：

```bash
docker compose --env-file .env.production -f docker-compose.prod.yaml up -d api
docker compose --env-file .env.production -f docker-compose.prod.yaml logs -f --tail=200 api
```

API 启动时会执行 GORM `AutoMigrate`。对于全新空库，它会创建当前模型表；对于已有生产库，仍应在发布前审核并显式执行版本化 SQL 迁移，不应把 `AutoMigrate` 当作完整的生产迁移方案。仓库中的迁移文件位于：

```text
gl-app/api/internal/models/migrations/
```

确认 API healthy 后启动前端与网关：

```bash
docker compose --env-file .env.production -f docker-compose.prod.yaml up -d app cms caddy
docker compose --env-file .env.production -f docker-compose.prod.yaml ps
```

观察日志：

```bash
docker compose --env-file .env.production -f docker-compose.prod.yaml logs -f --tail=200
```

按 `Ctrl+C` 只会退出日志跟随，不会停止容器。

### 9.4 首次验收

从服务器外部执行：

```bash
curl -fsSI https://app.example.com/
curl -fsSI https://cms.example.com/
curl -fsS 'https://app.example.com/app/v1/recommend/works?page=1&pageSize=1'
curl -fsS https://media.example.com/minio/health/live
```

浏览器验收至少覆盖：

1. App 首页加载和用户登录；
2. CMS 登录；
3. 图片上传、头像上传和图片展示；
4. 视频上传、Kafka 消费、FFmpeg 转码和 HLS 播放；
5. 点赞事件写入及刷新后状态一致；
6. MinIO 返回 URL 的域名为 `https://media.example.com`，不是 `localhost` 或 `minio`；
7. HTTPS 证书链正常，页面无 Mixed Content；
8. 重启单个容器后服务能自动恢复。

## 10. 日常发布流程

推荐每次发布按以下顺序：

### 10.1 拉取并确认版本

```bash
cd /opt/gulugulu/source
git fetch --all --tags --prune
git status --short
git checkout <已审核的-tag-或-commit>
git rev-parse --short=12 HEAD
```

工作区必须干净。不要在生产服务器上直接修改业务代码。

### 10.2 备份

涉及数据库结构、对象模型或重大业务逻辑的发布，必须先做 MySQL 和 MinIO 一致性备份。具体见“备份与恢复”。

### 10.3 构建新版本

将新的 Git 短提交号写入 `.env.production` 的 `IMAGE_TAG`：

```bash
cd /opt/gulugulu/source/gl-app/deploy
docker compose --env-file .env.production -f docker-compose.prod.yaml build --pull api app cms
```

### 10.4 更新服务

后端变更先更新 API，再更新前端：

```bash
docker compose --env-file .env.production -f docker-compose.prod.yaml up -d --no-deps api
docker compose --env-file .env.production -f docker-compose.prod.yaml ps api
docker compose --env-file .env.production -f docker-compose.prod.yaml logs --tail=200 api

docker compose --env-file .env.production -f docker-compose.prod.yaml up -d --no-deps app cms
docker compose --env-file .env.production -f docker-compose.prod.yaml ps app cms
```

`docker compose up -d` 不是严格意义的零停机部署。单机单副本 API 在容器替换期间可能有短暂连接失败。需要真正零停机时，应运行至少两个 API 副本、增加反向代理健康摘除，并使用滚动发布；这已超出简单单机 Compose 的可靠性边界。

### 10.5 发布后验证

```bash
docker compose --env-file .env.production -f docker-compose.prod.yaml ps
docker compose --env-file .env.production -f docker-compose.prod.yaml logs --since=10m api app cms caddy
curl -fsS 'https://app.example.com/app/v1/recommend/works?page=1&pageSize=1'
```

同时检查登录、上传、转码和关键业务指标。

## 11. 回滚

### 11.1 应用镜像回滚

前提是每个发布版本使用不可变 `IMAGE_TAG`，并保留旧镜像。

1. 将 `.env.production` 的 `IMAGE_TAG` 改回上一稳定版本；
2. 不执行 `build`；
3. 重新创建应用容器：

```bash
docker compose --env-file .env.production -f docker-compose.prod.yaml up -d --no-build --no-deps api app cms
docker compose --env-file .env.production -f docker-compose.prod.yaml ps
```

如果镜像来自私有仓库，先 `docker compose pull api app cms`。

### 11.2 数据库回滚

数据库迁移必须遵循向后兼容的 Expand/Contract 策略：

1. 先添加兼容字段或新表；
2. 发布同时兼容新旧结构的代码；
3. 数据回填；
4. 确认旧版本不再使用后，下一次发布再删除旧字段。

不要依赖“回滚镜像”自动回滚数据库。破坏性 DDL 只能通过经过审核的回滚 SQL 或备份恢复处理。

## 12. 备份与恢复

### 12.1 必须备份的内容

| 数据 | 是否必须 | 说明 |
| --- | --- | --- |
| MySQL | 必须 | 用户、作品、关系、审核状态等主数据 |
| MinIO | 必须 | 图片、视频、HLS、头像；必须与数据库保持一致 |
| `.env.production`、`prod/gulu.yaml` | 必须安全备份 | 建议加密保存到密码库/Secret Manager |
| Caddy `caddy_data` | 建议 | TLS 状态；丢失后可重新签发，但可能触发频率限制 |
| Redis | 视恢复目标 | 缓存可重建，但需评估未落库状态 |
| Kafka | 视恢复目标 | 队列可重建性取决于业务；有未消费任务时不能随意丢弃 |

### 12.2 MySQL 逻辑备份示例

```bash
cd /opt/gulugulu/source/gl-app/deploy
mkdir -p /opt/gulugulu/backup/mysql

docker compose --env-file .env.production -f docker-compose.prod.yaml exec -T mysql \
  sh -c 'exec mysqldump -uroot -p"$MYSQL_ROOT_PASSWORD" --single-transaction --routines --triggers --events --set-gtid-purged=OFF "$MYSQL_DATABASE"' \
  | gzip > "/opt/gulugulu/backup/mysql/gulugulu-$(date +%Y%m%d-%H%M%S).sql.gz"
```

定期检查备份不是空文件，并在隔离环境做恢复演练。大库应采用物理备份工具和 binlog，以实现更低 RPO 和时间点恢复。

### 12.3 MinIO 备份示例

不要只复制数据库而忽略对象存储。推荐使用 MinIO Client `mc mirror` 同步到独立服务器或云对象存储：

```bash
# 示例仅说明思路；生产凭据应从安全凭据文件或 Secret Manager 注入。
mc alias set source https://media.example.com <ACCESS_KEY> <SECRET_KEY>
mc alias set backup https://backup-s3.example.com <BACKUP_ACCESS_KEY> <BACKUP_SECRET_KEY>
mc mirror --overwrite --remove source/gulugulu-public backup/gulugulu-public
mc mirror --overwrite --remove source/gulugulu-media backup/gulugulu-media
mc mirror --overwrite --remove source/gulugulu-temp backup/gulugulu-temp
```

`--remove` 会删除备份端中源端不存在的对象，使用前必须确认是否符合保留策略。更稳妥的方案是开启对象版本控制、生命周期和异地复制。

### 12.4 恢复顺序

发生灾难恢复时：

1. 停止 Caddy、App、CMS 和 API，阻止新写入；
2. 恢复 MySQL；
3. 恢复 MinIO 对象；
4. 恢复必要配置；
5. 启动 MySQL、Redis、Kafka、MinIO；
6. 启动 API，检查迁移和 Worker 日志；
7. 启动 App、CMS、Caddy；
8. 验证数据库记录对应的对象都存在，再恢复流量。

停止应用层但保留基础设施：

```bash
docker compose --env-file .env.production -f docker-compose.prod.yaml stop caddy app cms api
```

不要使用 `docker compose down -v`。`-v` 会删除命名卷，可能造成不可恢复的数据丢失。

## 13. MinIO 控制台访问

Compose 只把控制台绑定到服务器的 `127.0.0.1:19001`。从管理员电脑建立 SSH 隧道：

```bash
ssh -L 19001:127.0.0.1:19001 deploy@server.example.com
```

然后在本机浏览器访问：

```text
http://127.0.0.1:19001
```

不要把 19001 加到公网安全组。MinIO 管理凭据应只交给运维人员，业务应用长期应使用权限受限的 Access Key，而不是 Root 用户。当前后端配置只有一组 MinIO 凭据字段，后续建议创建仅对三个业务桶有必要权限的服务账号并替换 Root 凭据。

## 14. 监控、日志与告警

单机 Compose 至少要监控：

- 宿主机 CPU、Load、内存、Swap；
- 系统盘和 Docker 数据盘使用率、inode、磁盘延迟；
- 容器重启次数和健康状态；
- API 请求量、P95/P99 延迟、4xx/5xx；
- MySQL 连接数、慢查询、锁等待、复制/备份状态；
- Redis 内存、命中率、淘汰和持久化失败；
- Kafka consumer lag、积压、磁盘；
- MinIO 容量、错误率、对象增长；
- FFmpeg 任务耗时、失败率和 `/tmp/gulugulu-media` 使用量；
- TLS 证书有效期和 Caddy 续期错误。

Docker 默认 `json-file` 日志可能无限增长。建议在 `/etc/docker/daemon.json` 设置日志轮转：

```json
{
  "log-driver": "json-file",
  "log-opts": {
    "max-size": "20m",
    "max-file": "5"
  }
}
```

修改 Docker daemon 配置会影响全机容器，应在维护窗口验证并重启 Docker：

```bash
sudo systemctl restart docker
```

已有容器通常需要重新创建后才应用新的日志选项。成熟环境建议将日志发送到 Loki/ELK/OpenSearch，并使用 Prometheus + Grafana 或云监控告警。

常用巡检命令：

```bash
docker compose --env-file .env.production -f docker-compose.prod.yaml ps
docker compose --env-file .env.production -f docker-compose.prod.yaml top
docker stats --no-stream
docker system df
df -h
```

## 15. 安全加固清单

上线前至少完成：

- [ ] SSH 禁止密码登录，禁止 root 远程登录，只允许密钥和受控来源 IP；
- [ ] 公网仅开放 80/443，数据库、Redis、Kafka、API、MinIO API/Console 不开放公网；
- [ ] 所有开发 Secret 已替换为独立随机值；
- [ ] `.env.production` 和 `prod/gulu.yaml` 权限为 600，且不在 Git 中；
- [ ] MySQL 使用业务账号，API 不使用 root；
- [ ] MinIO 后续改用最小权限服务账号，Root 账号仅用于管理；
- [ ] JWT 有效期合理，并设计刷新、撤销和密钥轮换机制；
- [ ] CMS 域名增加 VPN、IP 白名单、SSO 或额外身份认证；
- [ ] Caddy/前端限制最大请求体，与后端 100 MB 限制一致；
- [ ] Docker 和基础镜像定期升级，升级前做预发布验证；
- [ ] 镜像做漏洞扫描，依赖有定期更新机制；
- [ ] 备份加密、异地保存、限制访问，并做恢复演练；
- [ ] 管理操作和审核操作有审计日志；
- [ ] 服务器启用时间同步（systemd-timesyncd/chrony）；
- [ ] 对登录、验证码、上传、评论等接口配置限流和防滥用策略。

## 16. 常见故障排查

### 16.1 容器启动失败

```bash
docker compose --env-file .env.production -f docker-compose.prod.yaml ps -a
docker compose --env-file .env.production -f docker-compose.prod.yaml logs --tail=300 <服务名>
docker inspect gulugulu-<服务名>
```

Compose 生成的容器名不一定是 `gulugulu-<服务名>`，优先通过 `docker compose ps` 获取实际名称。

### 16.2 API 无法连接 MySQL/Redis/Kafka

确认 `gulu.yaml` 使用 Compose 服务名，而不是 `127.0.0.1`：

```text
mysql:3306
redis:6379
kafka:9092
```

在 API 容器内检查 DNS：

```bash
docker compose --env-file .env.production -f docker-compose.prod.yaml exec api getent hosts mysql redis kafka minio
```

检查各服务健康状态和日志。`depends_on` 只能帮助控制启动顺序，不能替代应用的连接重试和运行期监控。

### 16.3 图片 URL 出现 localhost 或 minio

检查 `prod/gulu.yaml`：

```yaml
Minio:
  Endpoint: media.example.com
  UseSSL: true
```

修改后重建 API 容器：

```bash
docker compose --env-file .env.production -f docker-compose.prod.yaml up -d --force-recreate --no-deps api
```

数据库中若保存的是对象 Key，新的响应会按新 Endpoint 生成 URL；如果数据库中保存了旧的完整 URL，则需要先备份再迁移数据。

### 16.4 预签名 URL 报 SignatureDoesNotMatch

常见原因：

- 签名时使用的 Host 与浏览器请求 Host 不一致；
- HTTP 被改成 HTTPS 或端口被改写；
- 代理未保留原始 Host；
- 服务器时间不准确；
- URL 被前端二次编码。

本文通过 `media.example.com -> Caddy -> MinIO` 保持内外 Host 一致。检查系统时间：

```bash
timedatectl status
```

### 16.5 Caddy 无法签发证书

检查：

```bash
docker compose --env-file .env.production -f docker-compose.prod.yaml logs --tail=300 caddy
getent hosts app.example.com cms.example.com media.example.com
sudo ss -lntup | grep -E ':(80|443)\b'
```

确认 DNS 指向本机、云安全组开放 80/443、没有其他 Nginx/Apache 占用端口。如果使用 CDN 代理，首次签证书时需确认 ACME 挑战能够正常到达 Caddy。

### 16.6 API 健康检查失败但进程存在

当前 API 健康检查会调用推荐作品接口，因此依赖 MySQL、Redis 等下游。依次检查：

```bash
docker compose --env-file .env.production -f docker-compose.prod.yaml exec api \
  curl -v 'http://127.0.0.1:8888/app/v1/recommend/works?page=1&pageSize=1'
docker compose --env-file .env.production -f docker-compose.prod.yaml logs --tail=300 api mysql redis kafka minio
```

长期建议在后端增加专门的 `/health/live` 和 `/health/ready`：liveness 只表示进程存活，readiness 检查关键依赖，避免业务接口变化影响健康判断。

### 16.7 视频一直处于处理中

Go API 与 Media Worker 在同一进程中。检查：

- Kafka 是否 healthy；
- `gulugulu-media-tasks` 是否有积压；
- API 日志中 FFmpeg/ffprobe 错误；
- `media_tmp` 空间和宿主机磁盘是否充足；
- MinIO Temp/Formal Bucket 是否可读写；
- `MediaQueue.Brokers` 是否为 `kafka:9092`；
- 容器 CPU 是否被限制过低。

验证镜像内工具：

```bash
docker compose --env-file .env.production -f docker-compose.prod.yaml exec api ffmpeg -version
docker compose --env-file .env.production -f docker-compose.prod.yaml exec api ffprobe -version
```

### 16.8 MySQL OOM 或性能异常

不要直接使用仓库开发 `my.cnf` 中的 `24G` Buffer Pool。查看实际配置和内存：

```bash
docker compose --env-file .env.production -f docker-compose.prod.yaml exec mysql \
  mysql -uroot -p -e "SHOW VARIABLES LIKE 'innodb_buffer_pool_size'; SHOW STATUS LIKE 'Threads_connected';"
docker stats --no-stream
```

依据服务器实际内存、连接数和查询负载调优，并开启慢查询分析。不要通过提高连接数掩盖慢 SQL。

## 17. 停止、重启与清理

安全重启单个服务：

```bash
docker compose --env-file .env.production -f docker-compose.prod.yaml restart api
```

停止应用但保留容器和数据：

```bash
docker compose --env-file .env.production -f docker-compose.prod.yaml stop
```

重新启动：

```bash
docker compose --env-file .env.production -f docker-compose.prod.yaml start
```

移除容器和网络但保留命名卷：

```bash
docker compose --env-file .env.production -f docker-compose.prod.yaml down
```

危险命令，生产环境通常禁止执行：

```bash
docker compose --env-file .env.production -f docker-compose.prod.yaml down -v
```

其中 `-v` 会删除 MySQL、MinIO、Kafka、Redis 和 Caddy 数据卷。除非已经确认是可销毁环境且完成备份，否则不要执行。

清理旧应用镜像前，先保留至少最近两个已验证版本：

```bash
docker image ls 'gulugulu-*'
```

不要在未核对磁盘和回滚需求时使用全局 `docker system prune -a --volumes`。

## 18. 从单机 Compose 向高可用演进

本方案适合早期生产、内部系统和可接受单机故障窗口的业务，但单机仍有明确上限：

- 服务器、Docker daemon、系统盘和公网 IP 都是单点；
- 单节点 MySQL、Kafka、MinIO 不具备真正高可用；
- API 同时承担 HTTP 和 Worker，无法独立扩缩容；
- Compose 不提供原生滚动发布、自动调度和跨节点故障转移。

业务增长后建议按顺序演进：

1. MySQL 使用云 RDS/主从集群，启用自动备份和 PITR；
2. Redis 使用托管版或 Sentinel/Cluster；
3. Kafka 使用至少 3 Broker 的托管服务或集群；
4. MinIO 使用分布式模式或云对象存储；
5. 将 API、Media Worker、Like Worker 拆成独立进程和镜像；
6. 应用部署迁移到 Kubernetes/ECS/Nomad 等编排平台；
7. 引入集中配置、Secret Manager、可观测性和自动化发布；
8. App/CMS 静态资源可进一步托管到 CDN/Object Storage。

## 19. 上线前最终核对表

- [ ] 三个域名 DNS 已生效，80/443 可从公网访问；
- [ ] 生产服务器为 amd64，或已完成 ARM64 镜像适配；
- [ ] `.env.production`、`prod/gulu.yaml` 未进入 Git 且权限正确；
- [ ] 所有开发密码、JWT Secret、Salt 已替换；
- [ ] `gulu.yaml` 的 MySQL、Redis、Kafka 地址使用 Compose 服务名；
- [ ] `Minio.Endpoint` 是媒体 HTTPS 公网域名且 `UseSSL: true`；
- [ ] 未挂载未经调优的开发 `my.cnf`；
- [ ] Compose 静默校验和 Caddyfile 校验通过；
- [ ] 后端、App、CMS 的测试和构建通过；
- [ ] MySQL、Redis、Kafka、MinIO、API 健康检查通过；
- [ ] 登录、上传、转码、播放、点赞、审核完整链路验收通过；
- [ ] 只有 22（受限）、80、443 对公网开放；
- [ ] 已配置日志轮转、主机监控、容器监控和告警；
- [ ] 已完成 MySQL + MinIO 联合备份；
- [ ] 已在隔离环境完成至少一次恢复演练；
- [ ] 已记录上一稳定 `IMAGE_TAG` 和明确的回滚步骤。

完成以上检查后，这套 Compose 部署可以作为 Gulugulu 的单机生产基线。任何涉及数据库结构、对象存储地址、域名、Secret 或持久卷的修改，都应先备份、在预发布环境验证，再进入生产。
