#!/bin/bash
# 開発者用：全OS向けにTagolyをクロスコンパイルするスクリプト

set -e

echo "🔨 Building tagoly binaries..."

# 既存の古いバイナリを削除（念のため）
rm -f tagoly-darwin-amd64 tagoly-darwin-arm64 tagoly-linux-amd64 tagoly-linux-arm64 tagoly-windows-amd64.exe

# Mac用
GOOS=darwin GOARCH=amd64 go build -o tagoly-darwin-amd64 ./cmd/tagoly
GOOS=darwin GOARCH=arm64 go build -o tagoly-darwin-arm64 ./cmd/tagoly

# Linux用
GOOS=linux GOARCH=amd64 go build -o tagoly-linux-amd64 ./cmd/tagoly
GOOS=linux GOARCH=arm64 go build -o tagoly-linux-arm64 ./cmd/tagoly

# Windows用
GOOS=windows GOARCH=amd64 go build -o tagoly-windows-amd64.exe ./cmd/tagoly

echo "✅ All builds completed successfully!"
ls -lh tagoly-*
