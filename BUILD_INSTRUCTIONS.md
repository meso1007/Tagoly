# Tagoly - Build & Install Instructions

## 前提条件
- Go 1.24.0 以上
- Git

## ビルド手順

### 1. ローカルでビルド（テスト用）
```bash
cd /path/to/tagoly
go build -o ./tagoly-local ./cmd/tagoly/
./tagoly-local tagdict
```

### 2. 開発版のインストール
```bash
cd /path/to/tagoly
go install ./cmd/tagoly/
tagoly tagdict
```

### 3. Homebrewで更新（リリース版）
```bash
brew upgrade meso1007/tagoly/tagoly
tagoly tagdict
```

## 使用方法

### tagdict - 対話的な検索（推奨）
```bash
tagoly tagdict
```
- "Commit Type" または "Scope" を選択
- タイプ/スコープの一覧から選択
- Enter で検索実行
- 矢印キーで結果をスクロール
- ESC で前の画面に戻る
- q で終了

### search - フラグでの検索
```bash
tagoly search -type feat
tagoly search -scope auth
tagoly search -subject login -limit 10
```

### コミット作成（デフォルト）
```bash
git add .
tagoly
```

## トラブルシューティング

### 古いバイナリが使用されている場合
Homebrewでインストールされたバイナリが優先されています。ローカルビルド版を使用してください：

```bash
# ローカルビルド版を使用
./tagoly-local tagdict

# または開発版をインストール
go install ./cmd/tagoly/
~/.gobin/tagoly tagdict
```

### サブコマンドが認識されない場合
Go モジュールを更新してリビルドしてください：

```bash
go mod tidy
go build -o ./tagoly-local ./cmd/tagoly/
```
