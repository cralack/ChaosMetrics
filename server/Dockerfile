FROM golang:1.19-alpine as builder
WORKDIR /app
COPY . /app
RUN go mod download
RUN go build -o chao main.go

FROM alpine:latest
RUN apk --no-cache add ca-certificates
CMD [ "./chao" ]