#!/bin/bash

# require
#apt install -y protobuf-compiler

# golang protoc plugin
#go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
#go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest

# go-micro plugin
#go install github.com/go-micro/generator/cmd/protoc-gen-micro@latest

# grpc-gateway plugin
#go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway@latest
#go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2@latest

# delete all pb go file
find . -type f -name "*pb*.go" -delete

# Loop through all .proto files in the current directory
for proto_file in *.proto
do
  if [ -f "$proto_file" ]; then
    # Generate Go code with protoc
    protoc -I "$GOPATH"/src -I . \
      --micro_out=./ \
      --go_out=. \
      --go-grpc_out=. \
      --grpc-gateway_out=logtostderr=true,register_func_suffix=Gw:. \
      "$proto_file"
  fi
done
