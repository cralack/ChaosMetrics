FROM golang:alpine AS builder

LABEL maintainer="cralack92@gmail.com"
WORKDIR /app
COPY . /app

RUN go env -w GO111MODULE=on \
    && go env -w GOPROXY=https://goproxy.cn,direct \
    && go env -w CGO_ENABLED=0 \
    && go env \
    && go mod tidy \
    && go build -o server .

RUN rm -f /app/api_key \
    && rm -f /app/email_key \
    && touch /app/api_key \
    && touch /app/email_key

FROM alpine:latest
RUN apk --no-cache add bash
WORKDIR /root/

COPY --from=builder /app/server ./server
COPY --from=builder /app/config_deploy.yaml ./config.yaml
COPY --from=builder /app/api_key ./
COPY --from=builder /app/email_key ./

ENTRYPOINT ["./server"]
CMD ["worker"]
