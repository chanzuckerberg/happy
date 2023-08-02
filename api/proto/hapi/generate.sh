#!/usr/bin/env bash

SRC_DIR=.
GO_DST_DIR=go/

rm -rf go
mkdir go

for proto in *.proto
do
    set -x
    protoc -I=$SRC_DIR --go_out=$GO_DST_DIR $SRC_DIR/$proto
    set +x
done
