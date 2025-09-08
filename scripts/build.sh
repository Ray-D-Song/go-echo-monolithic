#!/bin/bash

# Build script for go-echo-monolithic

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Configuration
PROJECT_NAME="go-echo-monolithic"
BUILD_DIR="bin"
SERVER_BINARY="server"
CLI_BINARY="cli"

echo -e "${GREEN}Building $PROJECT_NAME...${NC}"

# Create build directory
mkdir -p $BUILD_DIR

# Build server
echo -e "${YELLOW}Building server binary...${NC}"
go build -o $BUILD_DIR/$SERVER_BINARY ./cmd/server/

# Build CLI
echo -e "${YELLOW}Building CLI binary...${NC}"
go build -o $BUILD_DIR/$CLI_BINARY ./cmd/cli/

# Make binaries executable
chmod +x $BUILD_DIR/$SERVER_BINARY
chmod +x $BUILD_DIR/$CLI_BINARY

echo -e "${GREEN}Build completed successfully!${NC}"
echo -e "${GREEN}Binaries created:${NC}"
echo -e "  - $BUILD_DIR/$SERVER_BINARY"
echo -e "  - $BUILD_DIR/$CLI_BINARY"

# Show binary sizes
echo -e "${YELLOW}Binary sizes:${NC}"
ls -lh $BUILD_DIR/