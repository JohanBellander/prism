#!/bin/bash
set -e

# PRISM installer script
# Usage: curl -fsSL https://raw.githubusercontent.com/JohanBellander/prism/master/install.sh | bash

REPO="JohanBellander/prism"
BINARY_NAME="prism"

# Detect OS and architecture
OS=$(uname -s | tr '[:upper:]' '[:lower:]')
ARCH=$(uname -m)

case $ARCH in
    x86_64)
        ARCH="amd64"
        ;;
    aarch64|arm64)
        ARCH="arm64"
        ;;
    *)
        echo "Unsupported architecture: $ARCH"
        exit 1
        ;;
esac

# Determine install location
if [ -w "/usr/local/bin" ]; then
    INSTALL_DIR="/usr/local/bin"
elif [ -d "$HOME/.local/bin" ]; then
    INSTALL_DIR="$HOME/.local/bin"
else
    INSTALL_DIR="$HOME/bin"
    mkdir -p "$INSTALL_DIR"
fi

echo "Installing PRISM to $INSTALL_DIR..."

# Create temporary directory
TMP_DIR=$(mktemp -d)
cd "$TMP_DIR"

# Clone and build from source
echo "Cloning repository..."
git clone --depth 1 "https://github.com/$REPO.git" prism
cd prism

# Get version info
VERSION=$(git describe --tags --always 2>/dev/null || echo "dev")
COMMIT=$(git rev-parse --short HEAD 2>/dev/null || echo "none")
DATE=$(date -u +"%Y-%m-%dT%H:%M:%SZ")

echo "Building PRISM version $VERSION..."
if command -v go >/dev/null 2>&1; then
    go build -ldflags "-X main.version=$VERSION -X main.commit=$COMMIT -X main.date=$DATE" -o "$BINARY_NAME" ./cmd/prism
    mv "$BINARY_NAME" "$INSTALL_DIR/"
    chmod +x "$INSTALL_DIR/$BINARY_NAME"
else
    echo "Error: Go is required to build PRISM"
    echo "Install Go from https://go.dev/doc/install"
    exit 1
fi

# Cleanup
cd ~
rm -rf "$TMP_DIR"

echo ""
echo "✅ PRISM installed successfully to $INSTALL_DIR/$BINARY_NAME"
echo ""

# Check if install dir is in PATH
if ! echo "$PATH" | grep -q "$INSTALL_DIR"; then
    echo "⚠️  Add $INSTALL_DIR to your PATH:"
    echo ""
    if [ -f "$HOME/.bashrc" ]; then
        echo "  echo 'export PATH=\"$INSTALL_DIR:\$PATH\"' >> ~/.bashrc"
        echo "  source ~/.bashrc"
    elif [ -f "$HOME/.zshrc" ]; then
        echo "  echo 'export PATH=\"$INSTALL_DIR:\$PATH\"' >> ~/.zshrc"
        echo "  source ~/.zshrc"
    fi
    echo ""
fi

echo "Run 'prism --help' to get started!"
