FROM golang:alpine as builder

LABEL maintainer="cralack92@gmail.com"
WORKDIR /app
COPY . /app

RUN go env -w GO111MODULE=on \
    && go env -w GOPROXY=https://goproxy.cn,direct \
    && go env -w CGO_ENABLED=0 \
    && go env \
    && go mod tidy \
    && go build -o server . \

FROM alpine:latest
WORKDIR /root/
COPY --from=builder /app/chao ./
COPY --from=builder /app/config.yaml ./
CMD [ "./chao" ]