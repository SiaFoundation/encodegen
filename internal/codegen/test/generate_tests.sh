#!/bin/sh
set -eux
PACKAGE_DIR=$GOPATH/src/go.sia.tech/encodegen/internal/codegen

go test -test.v $PACKAGE_DIR