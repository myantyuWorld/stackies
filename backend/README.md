# Stackies Backend

Stackies プロジェクトのバックエンド API サーバー

## 技術スタック

- Go 1.21+
- Echo v4
- PostgreSQL (予定)
- Docker

## セットアップ

### ローカル開発環境

1. 依存関係のインストール

```bash
go mod download
```

2. サーバーの起動

```bash
go run main.go
```

### Docker を使用した開発

1. 開発環境の起動（ホットリロード対応）

```bash
docker-compose -f docker-compose.dev.yml up --build
```

2. 本番環境の起動

```bash
docker-compose up --build
```

サーバーは `http://localhost:8080` で起動します。

## API エンドポイント

- GET `/` - ウェルカムメッセージ

## 開発環境

- Go 1.21 以上
- Docker
- Docker Compose
- 推奨エディタ: VS Code + Go 拡張機能

## プロジェクト構造

```
backend/
├── main.go          # エントリーポイント
├── go.mod           # 依存関係管理
├── Dockerfile       # 本番環境用Dockerfile
├── Dockerfile.dev   # 開発環境用Dockerfile
├── docker-compose.yml       # 本番環境用Docker Compose設定
├── docker-compose.dev.yml   # 開発環境用Docker Compose設定
├── .air.toml        # ホットリロード設定
├── .gitignore       # Git除外設定
└── README.md        # プロジェクト説明
```
