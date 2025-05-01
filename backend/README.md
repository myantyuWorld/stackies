# Stackies Backend

Stackies プロジェクトのバックエンド API サーバー

## 技術スタック

- Go 1.21+
- Echo v4
- PostgreSQL 15
- Docker
- sql-migrate

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
make dev
```

2. 本番環境の起動

```bash
make up
```

### データベース操作

1. マイグレーションの実行

```bash
make migrate-up
```

2. マイグレーションのロールバック

```bash
make migrate-down
```

3. データベースコンテナに入る

```bash
make db-shell
```

4. データベースのリセット

```bash
make db-reset
```

サーバーは `http://localhost:8080` で起動します。

## API エンドポイント

- GET `/` - ウェルカムメッセージ

## 開発環境

- Go 1.21 以上
- Docker
- Docker Compose
- PostgreSQL 15
- 推奨エディタ: VS Code + Go 拡張機能

## プロジェクト構造

```
backend/
├── main.go          # エントリーポイント
├── config/          # 設定ファイル
│   └── database.go  # データベース設定
├── migrations/      # データベースマイグレーション
├── Dockerfile       # 本番環境用Dockerfile
├── Dockerfile.dev   # 開発環境用Dockerfile
├── docker-compose.yml       # 本番環境用Docker Compose設定
├── docker-compose.dev.yml   # 開発環境用Docker Compose設定
├── .air.toml        # ホットリロード設定
├── .gitignore       # Git除外設定
└── README.md        # プロジェクト説明
```

## データベース設計

### テーブル一覧

1. users

   - id
   - username
   - email
   - password_hash
   - created_at
   - updated_at

2. questions

   - id
   - user_id
   - title
   - content
   - created_at
   - updated_at

3. answers

   - id
   - question_id
   - user_id
   - content
   - created_at
   - updated_at

4. tags

   - id
   - name
   - created_at

5. question_tags
   - question_id
   - tag_id

## 開発ガイドライン

1. コード規約

   - Go の標準的なコーディング規約に従う
   - エラーハンドリングを適切に行う
   - テストを書く

2. コミットメッセージ

   - プレフィックスを使用（feat:, fix:, docs:, etc.）
   - 変更内容を明確に記述

3. プルリクエスト
   - 変更内容を詳細に記述
   - 関連する Issue を参照
   - レビューを依頼
