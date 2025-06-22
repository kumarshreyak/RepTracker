#!/bin/bash

# GymLog Backend - Protobuf Generation Script
# This script generates Go protobuf files from the proto definitions
# and places them in the correct backend directory structure.

set -e  # Exit on any error

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Script directory and project root
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
BACKEND_DIR="$(dirname "$SCRIPT_DIR")"

# Paths
PROTO_SOURCE_DIR="$BACKEND_DIR/proto"
PROTO_OUTPUT_DIR="$BACKEND_DIR/proto"

echo -e "${YELLOW}GymLog Backend - Protobuf Generation${NC}"
echo "=================================="
echo "Backend Dir: $BACKEND_DIR"
echo "Proto Source: $PROTO_SOURCE_DIR"
echo "Proto Output: $PROTO_OUTPUT_DIR"
echo ""

# Check if protoc is installed
if ! command -v protoc &> /dev/null; then
    echo -e "${RED}Error: protoc is not installed${NC}"
    echo "Please install Protocol Buffers compiler:"
    echo "  macOS: brew install protobuf"
    echo "  Ubuntu: sudo apt-get install protobuf-compiler"
    echo "  Or download from: https://github.com/protocolbuffers/protobuf/releases"
    exit 1
fi

# Check if Go protobuf plugins are installed
if ! command -v protoc-gen-go &> /dev/null; then
    echo -e "${RED}Error: protoc-gen-go is not installed${NC}"
    echo "Please install Go protobuf plugins:"
    echo "  go install google.golang.org/protobuf/cmd/protoc-gen-go@latest"
    echo "  go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest"
    exit 1
fi

if ! command -v protoc-gen-go-grpc &> /dev/null; then
    echo -e "${RED}Error: protoc-gen-go-grpc is not installed${NC}"
    echo "Please install Go protobuf plugins:"
    echo "  go install google.golang.org/protobuf/cmd/protoc-gen-go@latest"
    echo "  go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest"
    exit 1
fi

# Verify source proto file exists
PROTO_FILE="$PROTO_SOURCE_DIR/gymlog/v1/gymlog.proto"
if [ ! -f "$PROTO_FILE" ]; then
    echo -e "${RED}Error: Source proto file not found: $PROTO_FILE${NC}"
    exit 1
fi

echo -e "${YELLOW}Generating protobuf files...${NC}"

# Generate Go files from proto definitions
protoc \
    --proto_path="$PROTO_SOURCE_DIR" \
    --go_out=. \
    --go-grpc_out=. \
    --go_opt=module=gymlog-backend \
    --go-grpc_opt=module=gymlog-backend \
    gymlog/v1/gymlog.proto

if [ $? -eq 0 ]; then
    echo -e "${GREEN}✓ Protobuf files generated successfully${NC}"
    echo ""
    echo "Generated files:"
    find "$PROTO_OUTPUT_DIR" -name "*.pb.go" | sed 's|^|  |'
    echo ""
else
    echo -e "${RED}✗ Failed to generate protobuf files${NC}"
    exit 1
fi

# Verify the generated files exist
EXPECTED_FILES=(
    "$PROTO_OUTPUT_DIR/gymlog/v1/gymlog.pb.go"
    "$PROTO_OUTPUT_DIR/gymlog/v1/gymlog_grpc.pb.go"
)

echo -e "${YELLOW}Verifying generated files...${NC}"
for file in "${EXPECTED_FILES[@]}"; do
    if [ -f "$file" ]; then
        echo -e "${GREEN}✓ $(basename "$file")${NC}"
    else
        echo -e "${RED}✗ Missing: $(basename "$file")${NC}"
        exit 1
    fi
done

echo ""
echo -e "${GREEN}🎉 Protobuf generation completed successfully!${NC}"
echo ""
echo "Next steps:"
echo "  1. The generated files are ready to use in your Go code"
echo "  2. Import them using: pb \"gymlog-backend/proto/gymlog/v1\""
echo "  3. Run 'go build' to verify everything compiles correctly" 