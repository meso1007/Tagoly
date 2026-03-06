#!/bin/bash
# 開発者用：全OS向けにTagolyをクロスコンパイルするスクリプト

set -e

echo "🔨 Building tagoly binaries..."

# 既存の古いバイナリを削除（念のため）
rm -f tagoly-darwin-amd64 tagoly-darwin-arm64 tagoly-linux-amd64 tagoly-linux-arm64 tagoly-windows-amd64.exe

# Mac用
GOOS=darwin GOARCH=amd64 go build -o tagoly-darwin-amd64 cmd/tagoly/main.go
GOOS=darwin GOARCH=arm64 go build -o tagoly-darwin-arm64 cmd/tagoly/main.go

# Linux用
GOOS=linux GOARCH=amd64 go build -o tagoly-linux-amd64 cmd/tagoly/main.go
GOOS=linux GOARCH=arm64 go build -o tagoly-linux-arm64 cmd/tagoly/main.go

# Windows用
GOOS=windows GOARCH=amd64 go build -o tagoly-windows-amd64.exe cmd/tagoly/main.go

echo "✅ All builds completed successfully!"
ls -lh tagoly-*