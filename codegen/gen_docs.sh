#!/bin/bash
# Runs codegen tools
set -e
cd "$(dirname "$0")/.."

go run ./tools/docgen --version $(git describe --tags --abbrev=0) --output ./docs/plugins