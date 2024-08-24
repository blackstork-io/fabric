#!/bin/bash
# Runs protobuf linter
set -e
cd "$(dirname "${BASH_SOURCE[0]:-$0}")"
source ./setup.sh
cd ..

if is_ci; then
  buf lint --error-format=github-actions
else
  buf lint
fi