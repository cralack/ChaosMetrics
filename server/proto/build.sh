#!/bin/bash

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
