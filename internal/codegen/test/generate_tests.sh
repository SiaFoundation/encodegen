#!/bin/sh
set -eux
PACKAGE_DIR=$GOPATH/src/go.sia.tech/encodegen
TEST_DIR=$PACKAGE_DIR/internal/codegen/test

go test ../