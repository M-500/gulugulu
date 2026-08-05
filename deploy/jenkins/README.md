# Jenkins + Docker 发布说明

仓库内已经为三个项目分别提供独立流水线：

| 项目 | Jenkins Script Path | 默认容器名 | 默认宿主机端口 |
| --- | --- | --- | --- |
| App 前台 | `front/app/Jenkinsfile` | `gulugulu-app` | `8080` |
| CMS 后台 | `front/cms/Jenkinsfile` | `gulugulu-cms` | `8081` |
| Go 后端 | `gl-app/Jenkinsfile` | `gulugulu-api` | `8888` |

## Jenkins 前置条件

1. Jenkins Agent 已安装 Docker CLI，并有权限访问 Docker daemon。
2. 在 Jenkins Credentials 中创建“Username with password”类型的镜像仓库凭据，默认 ID 为 `docker-registry-credentials`。
3. 为三个项目分别创建 Pipeline/Multibranch Pipeline Job，并配置上表中的 Script Path。
4. 将 `REGISTRY` 改成实际镜像仓库域名，将 `IMAGE_NAMESPACE` 改成实际命名空间。
5. 需要自动部署时启用 `DEPLOY`。三个 Job 必须在能够访问同一个 Docker daemon 的 Agent 上部署。

如果 Jenkins 节点不能直接访问 Docker Hub，可通过 `NODE_IMAGE`、`NGINX_IMAGE`、`GO_IMAGE`、`RUNTIME_IMAGE` 参数填写公司 Harbor 或镜像代理中的完整基础镜像地址。

流水线镜像标签格式为：`分支名-构建号-Git短提交号`。默认还会推送 `latest`，可通过 `PUSH_LATEST` 关闭。

## 后端生产配置

后端镜像不会复制仓库里的开发配置，也不会把数据库、Redis、Kafka、MinIO 密码写入镜像。请在部署节点准备生产配置，例如：

```text
/opt/gulugulu/config/gulu.yaml
```

然后在后端 Job 的 `BACKEND_CONFIG_FILE` 参数中填写该绝对路径。流水线会只读挂载到容器内的 `/app/etc/gulu.yaml`。

生产配置中的媒体处理路径应使用镜像内命令：

```yaml
Media:
  FFmpegPath: ffmpeg
  FFprobePath: ffprobe
  TempDir: /tmp/gulugulu-media
```

后端运行镜像已经安装 `ffmpeg` 软件包，其中同时包含 `ffmpeg` 和 `ffprobe`。流水线会在推送镜像前执行版本检查。

## 容器网络与发布顺序

开启部署后，三个流水线会复用 `DOCKER_NETWORK` 参数指定的网络，默认是 `gulugulu`。推荐发布顺序：

1. 发布后端 `gulugulu-api`。
2. 确保 MySQL、Redis、Kafka 和 MinIO 可从该 Docker 网络访问，且生产配置使用可解析的容器名或实际地址。
3. 发布 App 和 CMS。

App/CMS 默认把 `/app/`、`/api/`、`/na/` 请求转发到 `http://gulugulu-api:8888`，把媒体桶路径转发到 `http://minio:9000`。如果基础设施名称不同，请修改流水线参数 `API_UPSTREAM` 和 `MINIO_UPSTREAM`。

`DEPLOY` 默认关闭，因此首次运行只会测试、构建和推送镜像，不会修改运行中的容器。
