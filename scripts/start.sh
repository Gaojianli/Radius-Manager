#!/bin/bash

echo "Starting RADIUS Management System..."

# Check if .env exists, if not copy from example
if [ ! -f ".env" ]; then
    echo "Creating .env file from example..."
    cp .env.example .env
    echo "Please edit .env file with your configuration before running again."
    exit 1
fi

# Build the application
echo "Building application..."
go build -o radius_mgnt .

if [ $? -eq 0 ]; then
    echo "Build successful. Starting server..."
    ./radius_mgnt
else
    echo "Build failed."
    exit 1
fi