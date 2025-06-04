#!/bin/bash

if ! command -v air &> /dev/null; then
    echo "Air is not installed. Installing now..."
    go install github.com/cosmtrek/air@latest

    if ! command -v air &> /dev/null; then
        echo "Failed to install Air. Please install manually with: go install github.com/cosmtrek/air@latest"
        exit 1
    fi

    echo "Air installed successfully!"
fi

mkdir -p tmp

# Set explicit dev flag for air to pass to the application
export AIR_SESSION_FLAG="-dev=true"
export GO_ENV="development"

echo "Starting application with Air for hot reloading in development mode..."
echo "Using flag: ${AIR_SESSION_FLAG}"
air
