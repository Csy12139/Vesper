#!/bin/bash

# Exit on error
set -xe

# Create build-test directory if it doesn't exist
mkdir -p build-test

echo "Compiling DN test file..."
go test -c -o build-test/dn ../DN
go test ./test
echo "Build complete! Binaries are in the build-test directory"