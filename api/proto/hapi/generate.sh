#!/usr/bin/env bash

SRC_DIR=.
GO_DST_DIR=go/
JS_DST_DIR=js/
TS_DST_DIR=ts/

rm -rf $GO_DST_DIR $JS_DST_DIR $TS_DST_DIR
mkdir  $GO_DST_DIR $JS_DST_DIR $TS_DST_DIR

for proto in *.proto
do
    set -x
    protoc -I=$SRC_DIR --go_out=$GO_DST_DIR $SRC_DIR/$proto
    protoc -I=$SRC_DIR --grpc-web_out=import_style=typescript,mode=grpcwebtext:$TS_DST_DIR $proto
    protoc -I=$SRC_DIR --js_out=import_style=commonjs:$JS_DST_DIR  --grpc-web_out=import_style=commonjs,mode=grpcwebtext:$JS_DST_DIR $proto
    set +x
done
