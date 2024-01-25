FROM golang:1.21-alpine as builder
LABEL maintainer="cralack92@gmail.com"
WORKDIR /app
COPY . /app
RUN go mod download
RUN go build -o chao main.go

FROM alpine:latest
WORKDIR /root/
COPY --from=builder /app/chao ./
COPY --from=builder /app/config.yaml ./
CMD [ "./chao" ]