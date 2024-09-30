#!/bin/bash

# Navigate to the project root directory (assumes this script is in the root)
cd "$(dirname "$0")"

# Build the Go project, specifying the output binary
go build -o tlp2-bin ./cmd/tlp2-bin

# Check if build succeeded and provide feedback
if [ $? -eq 0 ]; then
    echo "Build successful: tlp2-bin created."
else
    echo "Build failed."
fi

