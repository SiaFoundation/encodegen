#!/bin/sh
set -eux
go run ../../../ -s basic_struct/message.go -t Message -o basic_struct/message_buffer.go