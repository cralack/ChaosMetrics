FROM golang:alpine AS BUILDER

LABEL maintainer="cralack92@gmail.com"
WORKDIR /app
COPY . /app

RUN go env -w GO111MODULE=on \
    && go env -w GOPROXY=https://goproxy.cn,direct \
    && go env -w CGO_ENABLED=0 \
    && go env \
    && go mod tidy \
    && go build -o server .

RUN touch /app/api_key \
    && touch /app/email_key

FROM alpine:latest
RUN apk --no-cache add bash
WORKDIR /root/

COPY --from=BUILDER /app/server ./server
COPY --from=BUILDER /app/config.yaml ./config.yaml
COPY --from=BUILDER /app/api_key ./
COPY --from=BUILDER /app/email_key ./

CMD ["sh", "-c", "echo $API_KEY > /root/api_key && echo $EMAIL_KEY > /root/email_key && /bin/bash"]
