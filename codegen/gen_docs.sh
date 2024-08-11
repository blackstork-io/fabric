#!/bin/bash
# Runs codegen tools
set -e
cd "$(dirname "$0")/.."

go run ./tools/docgen --version $(git tag -l --sort=-creatordate | head -n1) --output ./docs/plugins