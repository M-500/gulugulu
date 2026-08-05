ARG NODE_IMAGE=node:20-alpine
ARG NGINX_IMAGE=nginx:1.27-alpine
ARG NPM_REGISTRY=https://registry.npmmirror.com

# App 与 CMS 分别在自己的构建阶段安装依赖，避免 node_modules 在阶段继承时丢失。
FROM ${NODE_IMAGE} AS app-build
WORKDIR /workspace/app
ENV NODE_ENV=development
ARG NPM_REGISTRY
COPY app/package.json app/package-lock.json ./
# npm 依赖安装和构建工具校验必须处于同一层，防止异常中断的 node_modules 被缓存。
RUN npm install --include=dev --registry="${NPM_REGISTRY}" 
ARG VITE_API_BASE_URL=""
ENV VITE_API_BASE_URL=${VITE_API_BASE_URL}
COPY app/ ./
RUN npm run build

FROM ${NODE_IMAGE} AS cms-build
WORKDIR /workspace/cms
ENV NODE_ENV=development
ARG NPM_REGISTRY
COPY cms/package.json cms/package-lock.json ./
RUN npm install --include=dev --registry="${NPM_REGISTRY}" 
ARG VUE_APP_API_BASE_URL=""
ARG VUE_APP_PUBLIC_PATH="/cms/"
ENV VUE_APP_API_BASE_URL=${VUE_APP_API_BASE_URL} \
    VUE_APP_PUBLIC_PATH=${VUE_APP_PUBLIC_PATH}
COPY cms/ ./
RUN npm run build

FROM ${NGINX_IMAGE} AS runtime

# 官方 Nginx 镜像启动时会将这里的变量写入配置模板。
ENV API_UPSTREAM=http://127.0.0.1:8888 \
    MINIO_UPSTREAM=http://127.0.0.1:9000 \
    NGINX_ENVSUBST_FILTER="^(API_UPSTREAM|MINIO_UPSTREAM)$"

COPY nginx.conf.template /etc/nginx/templates/default.conf.template
COPY --from=app-build /workspace/app/dist/ /usr/share/nginx/html/
COPY --from=cms-build /workspace/cms/dist/ /usr/share/nginx/html/cms/

EXPOSE 80

HEALTHCHECK --interval=30s --timeout=5s --start-period=10s --retries=3 \
  CMD wget -q -O /dev/null http://127.0.0.1/healthz || exit 1
