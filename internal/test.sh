#!/bin/sh

set -e

go install go.sia.tech/encodegen

rm -f test/encoding.go test/imported/encoding.go
go generate ./...
go test -v -fuzz=. -fuzztime=60s -bench=. -benchmem ./test -update
