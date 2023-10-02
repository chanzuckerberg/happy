#!/usr/bin/env bash

SRC_DIR=pkg/ent/proto/hapi
SHARED_DIR=../shared
GO_DST_DIR=../shared
TS_DST_DIR=hapi-proto/

# npm install

for proto in *.proto
do
    set -x
    gsed -i 's+github.com/chanzuckerberg/happy/api/pkg/ent/++g' $SRC_DIR/$proto
    protoc -I=$SRC_DIR --go_opt=M$SRC_DIR/$proto=hapi_protos --go_out=$SHARED_DIR  $SRC_DIR/$proto
    # protoc -I=$SRC_DIR --plugin=node_modules/ts-proto/protoc-gen-ts_proto --ts_proto_out=$TS_DST_DIR $SRC_DIR/$proto

    set +x
done
