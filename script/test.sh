#!/bin/bash

# Exit on error
set -xe

# Create build-test directory if it doesn't exist
mkdir -p build-test

echo "Compiling DN test file..."
go test -c -o build-test/dn ../DN
go test ./test
rm -rf test/logs && go test ./test/add_chunk_test.go
rm -rf test/logs && go test ./test/send_and_receive_file_test.go
bash ./script/sdkTest.sh
echo "Build complete! Binaries are in the build-test directory"