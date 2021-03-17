#!/bin/sh
set -eux
PACKAGE_DIR=$GOPATH/src/go.sia.tech/encodegen
TEST_DIR=$PACKAGE_DIR/internal/codegen/test

go run $PACKAGE_DIR -s $TEST_DIR/basic_struct -t Message,true -o $TEST_DIR/basic_struct/encoding.go
go run $PACKAGE_DIR -s $TEST_DIR/alias_struct -t Message,true -o $TEST_DIR/alias_struct/encoding.go
go run $PACKAGE_DIR -s $TEST_DIR/embedded_struct -t Message,true -o $TEST_DIR/embedded_struct/encoding.go
