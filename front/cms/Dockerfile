ARG NODE_IMAGE=node:20-alpine
ARG NGINX_IMAGE=nginx:1.27-alpine

FROM ${NODE_IMAGE} AS dependencies

WORKDIR /workspace
COPY package.json package-lock.json ./
RUN npm ci

FROM dependencies AS test

COPY . .
RUN npm run lint

FROM dependencies AS build

ARG VUE_APP_API_BASE_URL=""
ENV VUE_APP_API_BASE_URL=${VUE_APP_API_BASE_URL}
COPY . .
RUN npm run build

FROM ${NGINX_IMAGE} AS runtime

ENV API_UPSTREAM=http://gulugulu-api:8888 \
    MINIO_UPSTREAM=http://minio:9000 \
    NGINX_ENVSUBST_FILTER="^(API_UPSTREAM|MINIO_UPSTREAM)$"

COPY deploy/nginx.conf.template /etc/nginx/templates/default.conf.template
COPY --from=build /workspace/dist /usr/share/nginx/html

EXPOSE 80

HEALTHCHECK --interval=30s --timeout=5s --start-period=10s --retries=3 \
  CMD wget -q -O /dev/null http://127.0.0.1/ || exit 1
