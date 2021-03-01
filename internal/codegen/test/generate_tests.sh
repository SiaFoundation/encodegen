#!/bin/sh
set -eux
PACKAGE_DIR=$GOPATH/src/go.sia.tech/encodegen
TEST_DIR=$PACKAGE_DIR/internal/codegen/test

go run $PACKAGE_DIR -s $TEST_DIR/basic_struct/message.go -t Message -o $TEST_DIR/basic_struct/encoding.go
go run $PACKAGE_DIR -s $TEST_DIR/medium_struct/message.go -t Message -o $TEST_DIR/medium_struct/encoding.go
