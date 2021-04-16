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

echo "Instrumenting..."
cd fuzz
go get github.com/dvyukov/go-fuzz/go-fuzz-dep
go-fuzz-build

echo "Generating initial corpus..."
go get github.com/erggo/datafiller
go run gencorpus.go -out corpus

echo "Fuzzing..."
go-fuzz
