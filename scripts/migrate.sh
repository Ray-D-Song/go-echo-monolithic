#!/bin/bash

# Database migration script for go-echo-monolithic

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Configuration
CLI_BINARY="./bin/cli"
BUILD_SCRIPT="./scripts/build.sh"

# Function to check if CLI binary exists
check_cli_binary() {
    if [ ! -f "$CLI_BINARY" ]; then
        echo -e "${YELLOW}CLI binary not found. Building...${NC}"
        if [ -f "$BUILD_SCRIPT" ]; then
            bash "$BUILD_SCRIPT"
        else
            echo -e "${YELLOW}Build script not found. Building CLI manually...${NC}"
            mkdir -p bin
            go build -o "$CLI_BINARY" ./cmd/cli/
        fi
    fi
}

# Function to run migrations
migrate_up() {
    echo -e "${GREEN}Running database migrations...${NC}"
    $CLI_BINARY migrate
    echo -e "${GREEN}Migrations completed successfully!${NC}"
}

# Function to rollback migrations
migrate_down() {
    echo -e "${YELLOW}Rolling back database migrations...${NC}"
    echo -e "${RED}WARNING: This will drop all tables!${NC}"
    read -p "Are you sure you want to continue? (y/N): " -n 1 -r
    echo
    if [[ $REPLY =~ ^[Yy]$ ]]; then
        $CLI_BINARY rollback
        echo -e "${GREEN}Rollback completed successfully!${NC}"
    else
        echo -e "${YELLOW}Rollback cancelled.${NC}"
    fi
}

# Function to cleanup expired tokens
cleanup_tokens() {
    echo -e "${GREEN}Cleaning up expired tokens...${NC}"
    $CLI_BINARY cleanup
    echo -e "${GREEN}Token cleanup completed successfully!${NC}"
}

# Function to show usage
show_usage() {
    echo "Usage: $0 {up|down|cleanup|help}"
    echo ""
    echo "Commands:"
    echo "  up       - Run database migrations"
    echo "  down     - Rollback database migrations (drops all tables)"
    echo "  cleanup  - Clean up expired refresh tokens"
    echo "  help     - Show this help message"
}

# Main script logic
case "$1" in
    up)
        check_cli_binary
        migrate_up
        ;;
    down)
        check_cli_binary
        migrate_down
        ;;
    cleanup)
        check_cli_binary
        cleanup_tokens
        ;;
    help|--help|-h)
        show_usage
        ;;
    *)
        echo -e "${RED}Invalid command: $1${NC}"
        show_usage
        exit 1
        ;;
esac