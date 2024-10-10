#!/bin/bash

# Exit immediately if a command exits with a non-zero status
set -e
# Export all variables from the .env file to the current shell environment
set -a && source .env && set +a

# Database credentials
DB_NAME="testdb"
DB_PASSWORD=""
DB_HOST="localhost"
DB_PORT="5432"
PORT="3000"
DATABASE_URL="postgres://$DB_USER:$DB_PASSWORD@$DB_HOST:$DB_PORT/$DB_NAME?sslmode=disable"

# Function to drop the test database
cleanup() {
  echo "Shutting down the application..."
  [ -n "$SERVER_PID" ] && kill $SERVER_PID

  echo "Dropping test database..."
  psql -h "$DB_HOST" -d postgres -c "DROP DATABASE IF EXISTS $DB_NAME;"

  echo "Removing test binary..."
  rm -f bin/planner_server_test
}

# Register the cleanup function to be called on the EXIT signal
trap cleanup EXIT

# Create a new database for testing
echo "Creating test database..."
psql -h "$DB_HOST" -d postgres -c "CREATE DATABASE $DB_NAME;"

# Run migrations
echo "Running migrations..."
(
  cd sql/schema || exit
  goose postgres "$DATABASE_URL" up
)

# Build the Go script for populating the database
echo "Building the Go script..."
(
  cd scripts/populate_cuisines/
  go build -o bin/populate_cuisines populate_cuisines.go

  # Populate the database with cuisine data using the Go script
  echo "Populating database with cuisine data..."
  ./bin/populate_cuisines
)

# Build the test binary
echo "Building the test binary..."
go build -o bin/planner_server_test

# Run the application in the background
echo "Starting the application..."
bin/planner_server_test &
SERVER_PID=$!

# Give the server some time to start
sleep 1

# Run integration tests
echo "Running integration tests..."
hurl --test --jobs 1 --variable host=http://localhost:"$PORT" --variable email=testuser@testorg.com --variable password=verySafePassword1 --glob "tests/integration/**/*.hurl"
