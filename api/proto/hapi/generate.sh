#!/usr/bin/env bash

SRC_DIR=.
GO_DST_DIR=go/

rm -rf go grpcweb
mkdir go
mkdir grpcweb

for proto in *.proto
do
    set -x
    protoc -I=$SRC_DIR --go_out=$GO_DST_DIR $SRC_DIR/$proto
    protoc -I . --grpc-web_out=import_style=typescript,mode=grpcwebtext:grpcweb $proto
    set +x
done
