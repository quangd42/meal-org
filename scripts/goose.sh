#!/bin/bash

# Exit immediately if a command exits with a non-zero status
set -e
# Export all variables from the .env file to the current shell environment
set -a && source .env && set +a

go run github.com/pressly/goose/v3/cmd/goose@latest -dir "$SCHEMA_PATH" postgres "$DATABASE_URL" "$1"
