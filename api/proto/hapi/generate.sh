#!/usr/bin/env bash

SRC_DIR=.
GO_DST_DIR=go/
TS_DST_DIR=ts/

rm -rf $GO_DST_DIR $TS_DST_DIR
mkdir  $GO_DST_DIR $TS_DST_DIR

for proto in *.proto
do
    set -x
    protoc -I=$SRC_DIR --go_out=$GO_DST_DIR $SRC_DIR/$proto
    protoc -I=$SRC_DIR --plugin=node_modules/ts-proto/protoc-gen-ts_proto --ts_proto_out=$TS_DST_DIR $proto
    set +x
done
