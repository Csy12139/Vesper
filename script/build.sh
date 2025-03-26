#!/bin/bash

# Exit on error
set -xe

# Create build directory if it doesn't exist
mkdir -p build

# Generate protobuf code
echo "Generating protobuf code..."
#protoc --go_out=. --go_opt=paths=source_relative \
#    --go-grpc_out=. --go-grpc_opt=paths=source_relative \
#    ../proto/*.proto
protoc --proto_path=../proto --go_out=../proto --go_opt=paths=source_relative \
    --go-grpc_out=../proto --go-grpc_opt=paths=source_relative \
    ../proto/mn_service.proto

# Build DN binary
echo "Building DN..."
go build -o build/dn ../DN

# Build MN binary 
echo "Building MN..."
go build -o build/mn ../MN

echo "Build complete! Binaries are in the build directory"
