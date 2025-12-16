# ===============================
# 1. 基础构建阶段
# ===============================
FROM node:20-alpine AS base

ENV PNPM_HOME="/pnpm" \
    PNPM_REGISTRY_URL="https://registry.npmmirror.com" \
    PATH="$PNPM_HOME:$PATH" \
    CI=true

RUN corepack enable \
    && corepack prepare pnpm@10.25.0 --activate

WORKDIR /app

# 先只 copy package.json 和 lockfile
COPY package.json pnpm-lock.yaml* ./

# ===============================
# 2. 安装生产依赖
# ===============================
FROM base AS prod-deps
WORKDIR /app
RUN --mount=type=cache,id=pnpm,target=/pnpm/store \
    pnpm config set registry $PNPM_REGISTRY_URL \
    && pnpm install --prod

# ===============================
# 3. 构建阶段
# ===============================
FROM base AS build
WORKDIR /app

# 先 copy 生产依赖
COPY --from=prod-deps /app/node_modules /app/node_modules

# 再 copy其余源码
COPY . .

RUN --mount=type=cache,id=pnpm,target=/pnpm/store \
    pnpm config set registry $PNPM_REGISTRY_URL \
    && pnpm install \
    && pnpm run build

# ===============================
# 4. Nginx 生产镜像
# ===============================
FROM nginx:alpine
LABEL maintainer="cralack92@gmail.com"

RUN apk --no-cache add bash logrotate

# nginx 日志 rotate 配置
RUN echo "/var/log/nginx/*.log {\n    size 5000M\n    rotate 7\n    compress\n    missingok\n    notifempty\n    dateext\n}" > /etc/logrotate.d/nginx

# 复制 nginx 配置
COPY --from=build /app/nginx.conf /etc/nginx/conf.d/my.conf

# 复制构建产物
COPY --from=build /app/dist /usr/share/nginx/html
COPY --from=build /app/src/assets/tags /usr/share/nginx/html/src/assets/tags

# 复制脚本
COPY setupDataDragon.sh /usr/local/bin/setupDataDragon.sh
RUN chmod +x /usr/local/bin/setupDataDragon.sh

COPY docker-entrypoint.sh /usr/local/bin/
RUN chmod +x /usr/local/bin/docker-entrypoint.sh

WORKDIR /usr/share/nginx/html
ENTRYPOINT ["docker-entrypoint.sh"]