#!/bin/bash
# Runs codegen and docgen tools
set -e
cd "$(dirname "$0")"
./gen_code.sh
./gen_docs.sh
