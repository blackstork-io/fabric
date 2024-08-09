#!/bin/bash
# Runs codegen tools
set -e
cd "$(dirname "$0")/.."

if [ -z "$GOBIN" ]; then
  if [ -z "$GOPATH" ]; then
    if [ -z "$HOME" ]; then
      echo "HOME is unset"
      exit 1
    fi
    GOBIN="$HOME/go/bin"
  else
    GOBIN="$GOPATH/bin"
  fi
fi
PATH="$GOBIN:$PATH"

buf generate
mockery
go mod tidy