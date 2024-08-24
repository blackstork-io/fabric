#!/bin/bash
# Make codegen tools available
set -e
source "$(dirname "${BASH_SOURCE[0]:-$0}")/utils.sh"
# Codegen tools
install_tool mockery github.com/vektra/mockery/v2 "2.42.1"
install_tool buf github.com/bufbuild/buf/cmd/buf "1.32.2"
# Formatting tools
install_tool gci github.com/daixiang0/gci "0.13.4"
install_tool gofumpt mvdan.cc/gofumpt "0.6.0"
wait
