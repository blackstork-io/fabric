#!/bin/bash
# Installs tools
set -e
cd "$(dirname "$0")/.."
# Codegen tools
go install github.com/vektra/mockery/v2@v2.42.1 &
go install github.com/bufbuild/buf/cmd/buf@v1.32.2 &
# Formatting tools
go install github.com/daixiang0/gci@v0.13.4 &
go install mvdan.cc/gofumpt@v0.6.0 &
wait
