#!/bin/bash
# Formats the code
set -e
cd "$(dirname "${BASH_SOURCE[0]:-$0}")"
source ./setup.sh
cd ..

gofumpt -w .
gci write --skip-generated -s standard -s default -s "prefix(github.com/blackstork-io/fabric)" .

