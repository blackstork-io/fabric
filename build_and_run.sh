#!/bin/bash
set -e

cd "$( dirname "${BASH_SOURCE[0]}" )"
rm ./bin/* >/dev/null 2>&1 || true
go build -o ./bin/plugins ./cmd/plugins
go build -o ./bin/ .
./bin/fabric -path ./templates/ -plugins ./bin/plugins -document "test-document"
