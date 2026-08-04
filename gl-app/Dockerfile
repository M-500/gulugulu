ARG GO_IMAGE=golang:1.24-bookworm
ARG RUNTIME_IMAGE=ubuntu:24.04

FROM ${GO_IMAGE} AS dependencies

WORKDIR /workspace
COPY go.mod go.sum ./
ENV GO111MODULE=on \
    GOPROXY=https://goproxy.cn,https://proxy.golang.org,direct \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64
RUN go mod download

COPY . .
# RUN go test ./...

FROM dependencies AS build

ARG VERSION=dev
ARG GIT_COMMIT=unknown

RUN CGO_ENABLED=0 GOOS=linux go build \
    -trimpath \
    -ldflags "-s -w" \
    -o gulu-api ./api

FROM ${RUNTIME_IMAGE} AS runtime

ARG VERSION=dev
ARG GIT_COMMIT=unknown
ARG TZ=Asia/Shanghai

ENV TZ=${TZ}

# ffmpeg 软件包同时提供 ffmpeg 和 ffprobe，满足转码、HLS切片和媒体信息探测。
# 安装 tzdata 后将系统时区和 TZ 环境变量统一，确保 Go 日志及定时发布使用同一时区。
RUN apt-get update \
    && DEBIAN_FRONTEND=noninteractive apt-get install -y --no-install-recommends ca-certificates curl ffmpeg tzdata \
    && ln -snf "/usr/share/zoneinfo/${TZ}" /etc/localtime \
    && printf '%s\n' "${TZ}" > /etc/timezone \
    && rm -rf /var/lib/apt/lists/* \
    && groupadd --system --gid 10001 gulugulu \
    && useradd --system --uid 10001 --gid gulugulu --home-dir /app --no-log-init gulugulu \
    && mkdir -p /app/etc /tmp/gulugulu-media \
    && chown -R gulugulu:gulugulu /app /tmp/gulugulu-media

WORKDIR /app
COPY --from=build --chown=gulugulu:gulugulu gulu-api /app/gulu-api

LABEL org.opencontainers.image.version=${VERSION} \
      org.opencontainers.image.revision=${GIT_COMMIT}

USER gulugulu
EXPOSE 8888

HEALTHCHECK --interval=30s --timeout=5s --start-period=20s --retries=3 \
  CMD curl -fsS "http://127.0.0.1:8888/app/v1/recommend/works?page=1&pageSize=1" >/dev/null || exit 1

ENTRYPOINT ["/app/gulu-api"]
CMD ["-f", "/app/etc/gulu.yaml"]
