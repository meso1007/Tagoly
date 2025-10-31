# Tagoly

Tagoly は、Git コミットをスマートに支援する CLI ツールです。  
スコープ検出・カスタムタグ対応・インタラクティブな選択など、手動でのコミットメッセージ作成を効率化します。

## 🚀 主な機能

- **自動スコープ検出**  
  変更されたファイルパスから自動的にスコープを判定

- **カスタムタグ対応**  
  `.tagolyrc` に自分専用のタグを定義可能（例：ci, perf など）

- **対話的コミット生成**  
  commit type / scope / message を順にインタラクティブに選択

- **スマートスコープ選択**  
  複数スコープが含まれる場合は最も多いものを自動選択  
  必要に応じて手動で選択可能

## インストール方法

### **MacOS**

#### 1. Homebrew
```bash
# Tap を追加
brew tap meso1007/tagoly

# インストール
brew install meso1007/tagoly/tagoly

# 動作確認
tagoly --version
tagoly
```

#### 2. 手動インストール
##### Apple Silicon (M1/M2)
```bash
mv tagoly-darwin-arm64 /usr/local/bin/tagoly && chmod +x /usr/local/bin/tagoly
```
##### Intel
```bash
mv tagoly-darwin-amd64 /usr/local/bin/tagoly && chmod +x /usr/local/bin/tagoly
```
--------

### **Linux**
#### 1. Homebrew (Linuxbrew)
```bash
brew tap meso1007/tagoly
brew install meso1007/tagoly/tagoly
```
#### 2. 手動インストール
```bash
mv tagoly-linux-amd64 /usr/local/bin/tagoly && chmod +x /usr/local/bin/tagoly
```
--------

### **Windows**
#### 1. Scoop
```powershell
# バケットを追加
scoop bucket add tagoly https://github.com/meso1007/scoop-tagoly

# インストール
scoop install tagoly/tagoly

# 動作確認
tagoly --version
tagoly
```
#### 2. 手動インストール
```powershell
Move-Item .\tagoly.exe "C:\Program Files\tagoly\tagoly.exe"
```

--------

## 設定ファイル .tagolyrc
```json
{
  "customTags": ["ci", "perf"]
}
```

## 使用方法
```bash
git add .
tagoly
```

## リポジトリ

- Tagoly 本体: [https://github.com/meso1007/tagoly](https://github.com/meso1007/tagoly)  
- Homebrew Tap: [https://github.com/meso1007/homebrew-tagoly](https://github.com/meso1007/homebrew-tagoly)  
- Scoop バケット: [https://github.com/meso1007/scoop-tagoly](https://github.com/meso1007/scoop-tagoly)

---
