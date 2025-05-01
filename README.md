# Stackies

業務経歴を可視化・共有するための Web アプリケーション

## プロジェクト概要

### 目的

- エンジニアの業務経歴を構造化して管理
- 技術スタックやプロジェクト経験を可視化
- キャリアパスの共有と分析

### 主な機能

- 業務経歴の登録・編集
- 技術スタックの可視化
- プロジェクト経験の詳細管理
- スキルセットの分析
- キャリアパスの共有

## 技術スタック

### フロントエンド

- Vue 3
- TypeScript
- Vite
- Tailwind CSS
- Chart.js（スキル可視化用）

### バックエンド

- Go
- Echo
- PostgreSQL
- Docker

## 開発環境のセットアップ

### 前提条件

- Node.js 18+
- Go 1.21+
- Docker
- Docker Compose

### フロントエンド

```bash
cd frontend
npm install
npm run dev
```

### バックエンド

```bash
cd backend
make dev
```

## プロジェクト構造

```
stackies/
├── frontend/          # Vue 3フロントエンド
│   ├── src/          # ソースコード
│   │   ├── components/  # UIコンポーネント
│   │   ├── views/      # ページコンポーネント
│   │   ├── stores/     # 状態管理
│   │   └── types/      # TypeScript型定義
│   ├── public/       # 静的ファイル
│   └── package.json  # 依存関係
│
└── backend/          # Goバックエンド
    ├── src/         # ソースコード
    │   ├── handlers/  # APIハンドラ
    │   ├── models/    # データモデル
    │   └── services/  # ビジネスロジック
    ├── migrations/  # データベースマイグレーション
    └── Dockerfile   # コンテナ設定
```

## データベース

### 開発環境

- PostgreSQL 15
- ポート: 5432
- データベース名: stackies_dev
- ユーザー: postgres
- パスワード: postgres

### マイグレーション

```bash
cd backend
make migrate-up    # マイグレーション実行
make migrate-down  # マイグレーションロールバック
```

## デプロイ

### 本番環境

- AWS ECS
- AWS RDS
- AWS S3

### CI/CD

- GitHub Actions
- 自動テスト
- 自動デプロイ

## コントリビューション

1. このリポジトリをフォーク
2. 新しいブランチを作成 (`git checkout -b feature/amazing-feature`)
3. 変更をコミット (`git commit -m 'Add some amazing feature'`)
4. ブランチにプッシュ (`git push origin feature/amazing-feature`)
5. プルリクエストを作成

## ライセンス

MIT License
