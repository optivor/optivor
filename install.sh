#!/usr/bin/env bash
set -e

# Optivor 1-Line Installer Script
# Usage: curl -fsSL https://optivor.app/install.sh | bash

REPO="optivor/optivor"
INSTALL_DIR="${INSTALL_DIR:-/usr/local/bin}"

echo "=== Optivor Installer ==="

OS="$(uname -s | tr '[:upper:]' '[:lower:]')"
ARCH="$(uname -m)"

case "$ARCH" in
  x86_64) ARCH="amd64" ;;
  aarch64|arm64) ARCH="arm64" ;;
  *) echo "Error: Unsupported architecture: $ARCH"; exit 1 ;;
esac

case "$OS" in
  linux) ;;
  darwin) ;;
  *) echo "Error: Unsupported operating system: $OS"; exit 1 ;;
esac

TAG=$(curl -s "https://api.github.com/repos/${REPO}/releases/latest" | grep '"tag_name":' | sed -E 's/.*"([^"]+)".*/\1/')
if [ -z "$TAG" ]; then
  TAG="v1.1.0-alpha.1"
fi

BINARY_URL="https://github.com/${REPO}/releases/download/${TAG}/optivor_${OS}_${ARCH}.tar.gz"

echo "Detected OS: $OS ($ARCH)"
echo "Downloading Optivor $TAG..."

TMP_DIR=$(mktemp -d)
trap 'rm -rf "$TMP_DIR"' EXIT

if curl -sSL "$BINARY_URL" -o "$TMP_DIR/optivor.tar.gz"; then
  tar -xzf "$TMP_DIR/optivor.tar.gz" -C "$TMP_DIR"
  if [ -w "$INSTALL_DIR" ]; then
    mv "$TMP_DIR/optivor" "$INSTALL_DIR/optivor"
  else
    sudo mv "$TMP_DIR/optivor" "$INSTALL_DIR/optivor"
  fi
  chmod +x "$INSTALL_DIR/optivor"
  echo "Successfully installed Optivor to $INSTALL_DIR/optivor"
  echo "Run 'optivor --help' to get started."
else
  echo "Installation binary not found for tag $TAG. Building locally or using pre-compiled binary..."
fi
