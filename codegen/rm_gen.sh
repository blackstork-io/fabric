#!/bin/bash
# Removes all generated files

cd "$(dirname "$0")/.."

grep -R --with-filename --files-with-matches --no-messages --include "*.go" -E -e '^// Code generated .* DO NOT EDIT.$' . | xargs rm
find ./mocks -type d -empty -exec rmdir {} +

find ./docs/plugins -mindepth 1 \( -maxdepth 1 -type d -o -not -name "*.md" \) -exec rm -rf {} +
