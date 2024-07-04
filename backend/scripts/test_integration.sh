#!/bin/bash

# Exit immediately if a command exits with a non-zero status
set -e

# Database credentials
DB_NAME="testdb"
DB_USER="quang-dang"
DB_PASSWORD=""
DB_HOST="localhost"
DB_PORT="5432"
PORT="3000"

# Function to drop the test database
cleanup() {
  # Kill the application
  if [ -n "$SERVER_PID" ]; then
    echo "Shutting down the application..."
    kill $SERVER_PID
  fi

  # Drop the database
  echo "Dropping test database..."
  psql -h "$DB_HOST" -d postgres -c "DROP DATABASE IF EXISTS $DB_NAME;"

  # Remove the test binary
  rm bin/planner_server_test
}

# Register the cleanup function to be called on the EXIT signal
trap cleanup EXIT

# Create a new database for testing
echo "Creating test database..."
psql -h "$DB_HOST" -d postgres -c "CREATE DATABASE $DB_NAME;"

# Run migrations
echo "Running migrations..."
cd sql/schema || exit
goose postgres "postgres://$DB_USER:$DB_PASSWORD@$DB_HOST:$DB_PORT/$DB_NAME" up
cd ../../

# Build the test binary
go build -o bin/planner_server_test

# Run the application in the background
echo "Starting the application..."
PORT=$PORT DATABASE_URL="postgres://$DB_USER:$DB_PASSWORD@$DB_HOST:$DB_PORT/$DB_NAME" bin/planner_server_test &
SERVER_PID=$!

# Give the server some time to start
sleep 1

# Run integration tests
echo "Running integration tests..."
hurl --test --variable host=http://localhost:"$PORT" --variable username=jbergey2 --variable password=verySafePassword1 --glob "integration-tests/**/*.hurl"
