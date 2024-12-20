# 基础构建阶段
FROM node:20-alpine AS base
ENV PNPM_HOME="/pnpm"
ENV PATH="$PNPM_HOME:$PATH"
RUN corepack enable
COPY . /app
WORKDIR /app


FROM base AS prod-deps
RUN --mount=type=cache,id=pnpm,target=/pnpm/store pnpm install --prod

FROM base AS build
COPY --from=prod-deps /app/node_modules /app/node_modules
RUN --mount=type=cache,id=pnpm,target=/pnpm/store pnpm install && pnpm run build


FROM nginx:alpine
LABEL maintainer="cralack92@gmail.com"
COPY --from=base   /app/nginx.conf /etc/nginx/conf.d/my.conf
COPY --from=build  /app/dist /usr/share/nginx/html
COPY --from=build  /app/src/assets/tags /usr/share/nginx/html/src/assets/tags

RUN apk --no-cache add bash
RUN ls -al /usr/share/nginx/html

COPY setupDataDragon.sh /usr/local/bin/setupDataDragon.sh
RUN chmod +x /usr/local/bin/setupDataDragon.sh

COPY docker-entrypoint.sh /usr/local/bin/
RUN chmod +x /usr/local/bin/docker-entrypoint.sh

WORKDIR /usr/share/nginx/html
ENTRYPOINT ["docker-entrypoint.sh"]