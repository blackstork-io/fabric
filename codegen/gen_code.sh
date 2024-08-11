#!/bin/bash
# Runs codegen tools
set -e
cd "$(dirname "${BASH_SOURCE[0]:-$0}")"
source ./setup.sh
cd ..

buf generate
mockery
go mod tidy