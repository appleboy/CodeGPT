#!/bin/bash

REPO="appleboy/CodeGPT"
API_URL="https://api.github.com/repos/$REPO/releases/latest"

OS=$(uname -s)
ARCH=$(uname -m)

SELECTED_BIN=""
TARGET_FILE="codegpt"
DEST_DIR="/usr/bin"

case "$OS" in
"FreeBSD")
    if [ "$ARCH" = "amd64" ]; then
        PATTERN="freebsd-amd64"
    fi
    ;;
"Linux")
    case "$ARCH" in
    "x86_64")
        PATTERN="linux-amd64"
        ;;
    "aarch64")
        PATTERN="linux-arm64"
        ;;
    "armv7l")
        PATTERN="linux-arm-7"
        ;;
    "armv6l")
        PATTERN="linux-arm-6"
        ;;
    "armv5tel")
        PATTERN="linux-arm-5"
        ;;
    esac
    ;;
esac

if [ -z "$PATTERN" ]; then
    echo "Unsupported OS or architecture: $OS $ARCH"
    exit 1
fi

echo "Fetching latest release from GitHub..."
ASSET_URL=$(curl -s "$API_URL" | grep "browser_download_url" | grep "$PATTERN" | grep -v '\.xz' | cut -d '"' -f 4)

if [[ -z "$ASSET_URL" ]]; then
    echo "Failed to find a compatible binary for $PATTERN."
    exit 1
fi

echo "Downloading binary from: $ASSET_URL"

TEMP_BIN="/tmp/$TARGET_FILE"
curl -L --retry 3 -o "$TEMP_BIN" "$ASSET_URL"

if [ $? -ne 0 ]; then
    echo "Download failed."
    exit 1
fi

chmod 755 "$TEMP_BIN"
sudo mv "$TEMP_BIN" "$DEST_DIR/$TARGET_FILE"

if [ $? -eq 0 ]; then
    echo "Installation successful! You can now run CodeGPT with '$TARGET_FILE'."
    $TARGET_FILE version
else
    echo "Installation failed."
    exit 1
fi
