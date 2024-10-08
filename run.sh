#!/bin/bash

# Set the temporary directory
tmp_dir="dist"

# Create the temporary directory if it doesn't exist
mkdir -p $tmp_dir

# Build the Go application (without .exe)
echo "Building the application..."
go build -o "$tmp_dir/main" ./src

# Check if the build was successful
if [ $? -eq 0 ]; then
    echo "Build successful. Starting the application..."
    # Start the application
    ./"$tmp_dir/main"
else
    echo "Build failed. Please check the errors."
fi