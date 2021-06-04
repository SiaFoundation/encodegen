#!/bin/sh

set -e

trap cleanup 1 2 3 6

cleanup() {
  rm -f test/encoding.go test/imported/encoding.go
  rm -f ../test/encoding.go ../test/imported/encoding.go
  exit 1
}

go install go.sia.tech/encodegen

echo "Generating methods..."
rm -f test/encoding.go test/imported/encoding.go
go generate ./...

echo "Fuzzing..."
# install gotip with: go get golang.org/dl/gotip && gotip download dev.fuzz
gotip test -v -fuzz=FuzzUnmarshalSimple ./test/.