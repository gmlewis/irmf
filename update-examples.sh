#!/bin/bash -ex
go run cmd/update-examples/main.go
pt -l '##' examples/*/README.md | sort | sed -e 's|/README.md||'
