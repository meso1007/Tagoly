#!/bin/bash

# Tagoly Installer Script
# ユーザーはこれを実行するだけでOK

set -e

REPO="meso1007/tagoly"
BINARY_NAME="tagoly"

# 1. OSとアーキテクチャの検出
OS="$(uname -s)"
ARCH="$(uname -m)"

case "$OS" in
    Linux)  PLATFORM="linux" ;;
    Darwin) PLATFORM="darwin" ;;
    MINGW*|MSYS*|CYGWIN*) PLATFORM="windows" ;;
    *) echo "Unsupported OS: $OS"; exit 1 ;;
esac

case "$ARCH" in
    x86_64)  ARCH_TAG="amd64" ;;
    arm64|aarch64) ARCH_TAG="arm64" ;;
    *) echo "Unsupported Architecture: $ARCH"; exit 1 ;;
esac

# Windowsの場合は拡張子をつける
if [ "$PLATFORM" = "windows" ]; then
    BINARY_NAME="${BINARY_NAME}.exe"
fi

echo "🚀 Installing Tagoly for $PLATFORM/$ARCH_TAG..."

# 2. 最新のリリースURLを取得（GitHub API）
LATEST_URL=$(curl -s "https://api.github.com/repos/$REPO/releases/latest" | \
    grep "browser_download_url" | \
    grep "$PLATFORM" | \
    grep "$ARCH_TAG" | \
    cut -d '"' -f 4)

if [ -z "$LATEST_URL" ]; then
    echo "❌ Error: Could not find a release for your platform."
    exit 1
fi

echo "Found latest release: $LATEST_URL"

# 3. ダウンロードとインストール
TEMP_DIR=$(mktemp -d)
# Tarballではなくバイナリを直接ダウンロードする
curl -L "$LATEST_URL" -o "$TEMP_DIR/$BINARY_NAME"

# 実行権限を付与
chmod +x "$TEMP_DIR/$BINARY_NAME"

# インストール先（Mac/Linuxなら /usr/local/bin, Windowsならカレント等）
INSTALL_DIR="/usr/local/bin"

echo "Installing to $INSTALL_DIR..."

# 管理者権限が必要かチェック
if [ -w "$INSTALL_DIR" ]; then
    mv "$TEMP_DIR/$BINARY_NAME" "$INSTALL_DIR/"
else
    echo "🔑 Sudo permission required to move binary to $INSTALL_DIR"
    sudo mv "$TEMP_DIR/$BINARY_NAME" "$INSTALL_DIR/"
fi

# 4. 完了メッセージ
echo ""
echo "✅ Tagoly has been installed successfully!"
echo "Try running: tagoly --version"

# クリーンアップ
rm -rf "$TEMP_DIR"