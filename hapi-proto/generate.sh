#!/usr/bin/env bash
set -euox pipefail

SRC_DIR=hapi
# GO_DST_DIR=dist/go/
# TS_DST_DIR=hapi-proto/

# rm -rf $GO_DST_DIR $TS_DST_DIR
# mkdir -p $GO_DST_DIR $TS_DST_DIR

npm install
go install github.com/golang/protobuf/protoc-gen-go@master
go install github.com/infobloxopen/protoc-gen-gorm@main
GOPATH=$(go env GOPATH)
PROTOGEN_GO_VERSION="$(ls $GOPATH/pkg/mod/github.com/infobloxopen | grep -m 2 protoc-gen-go | tail -1)"
PROTOGEN_GO_PATH="$GOPATH/pkg/mod/github.com/infobloxopen/$PROTOGEN_GO_VERSION/proto"
MODEL_PROTO_PATH="hapi/service-event.proto"

# protoc -I="${PROTOGEN_GO_PATH}" -I="." --go_out="." --gorm_out="engine=postgres:." ${MODEL_PROTO_PATH}

for proto in $(find ${SRC_DIR} -name '*.proto');
do
    protoc -I="${PROTOGEN_GO_PATH}" -I="." --go_out="." --gorm_out="engine=postgres:." ${proto}
done

# GORM_STRUCT_DIR="../api/pkg/model/"
# mkdir -p $GORM_STRUCT_DIR
# mv hapi-protos/*.pb.gorm.go $GORM_STRUCT_DIR
