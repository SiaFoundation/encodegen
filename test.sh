#!/bin/sh

set -e

go install .

echo "Reflection:"
rm -f test/encoding.go test/imported/encoding.go
go test -bench=. -benchmem ./test

echo "Codegen:"
go generate ./...
go test -bench=. -benchmem ./test
rm -f test/encoding.go test/imported/encoding.go
