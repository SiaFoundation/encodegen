#!/bin/sh

set -e

trap cleanup 1 2 3 6

cleanup() {
  rm -f test/encoding.go test/imported/encoding.go
  rm -f ../test/encoding.go ../test/imported/encoding.go
  exit 1
}

go install .

echo "Generating methods..."
rm -f test/encoding.go test/imported/encoding.go
go generate ./...

echo "Instrumenting..."
cd fuzz
go-fuzz-build

echo "Generating initial corpus..."
go run gencorpus.go -out corpus

echo "Fuzzing..."
go-fuzz
