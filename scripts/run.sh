#!/bin/sh
set -eux
gofmt -s -w .
reset
go run cmd/encodegen/main.go "${@}"