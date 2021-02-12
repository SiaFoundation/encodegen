#!/bin/sh
set -eux
gofmt -s -w .
reset
go build -o encodegen cmd/encodegen/main.go "${@}"