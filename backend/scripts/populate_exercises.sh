#!/bin/bash

echo "Starting exercise population script..."

# Change to the backend directory
cd "$(dirname "$0")/.."

# Check if the JSON file exists
if [ ! -f "scripts/strengthlog_exercises.json" ]; then
    echo "Error: strengthlog_exercises.json not found in scripts directory"
    exit 1
fi

# Build and run the population script
echo "Building population script..."
go build -o scripts/populate_exercises scripts/populate_exercises.go

if [ $? -ne 0 ]; then
    echo "Error: Failed to build population script"
    exit 1
fi

echo "Running population script..."
./scripts/populate_exercises

# Clean up the binary
rm -f scripts/populate_exercises

echo "Exercise population completed!" 