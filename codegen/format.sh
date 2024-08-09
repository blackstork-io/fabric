#!/bin/bash
# Formats the code
set -e
cd "$(dirname "$0")/.."
gofumpt -w .
gci write --skip-generated -s standard -s default -s "prefix(github.com/blackstork-io/fabric)" .

