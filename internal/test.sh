#!/bin/sh

set -e

go install go.sia.tech/encodegen

echo "Reflection:"
rm -f test/encoding.go test/imported/encoding.go
go test -v -bench=. -benchmem ./test -update

echo "Codegen:"
go generate ./...
go test -v -bench=. -benchmem ./test
rm -f test/encoding.go test/imported/encoding.go
