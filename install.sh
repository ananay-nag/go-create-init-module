#!/bin/bash

set -e

CLI_NAME="go-set-mod"
REPO_URL="https://github.com/ananay-nag/go-create-init-module.git"
INSTALL_DIR="$(dirname $(which go))"  # Get Go’s installation directory

echo "Installing $CLI_NAME..."

# Ensure dependencies are installed
if ! command -v git &> /dev/null; then
    echo "Error: git is not installed. Please install git first."
    exit 1
fi

if ! command -v go &> /dev/null; then
    echo "Error: Go is not installed. Please install Go first."
    exit 1
fi

# Create a temp directory
TEMP_DIR=$(mktemp -d)
cd "$TEMP_DIR"

# Clone the repository
echo "Cloning repository..."
git clone "$REPO_URL" "$CLI_NAME"
cd "$CLI_NAME"

# Build the Go project
echo "Building $CLI_NAME..."
go build -o $CLI_NAME main.go

# Move binary to Go's installation directory
echo "Moving $CLI_NAME to $INSTALL_DIR..."
sudo mv $CLI_NAME "$INSTALL_DIR/"

# Cleanup
cd ..
rm -rf "$TEMP_DIR"

echo "Installation complete! Run '$CLI_NAME <module-name>' to use."
