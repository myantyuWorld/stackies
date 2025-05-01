.PHONY: dev up down build clean db-shell db-reset

# 開発環境の起動
dev:
	docker-compose -f docker-compose.dev.yml up --build

# 本番環境の起動
up:
	docker-compose up --build

# コンテナの停止
down:
	docker-compose down

# イメージのビルド
build:
	docker-compose build

# 開発環境のクリーンアップ
clean:
	docker-compose down -v
	rm -rf tmp/
	go clean -cache -testcache -modcache

# 開発環境のログ表示
logs:
	docker-compose -f docker-compose.dev.yml logs -f

# 本番環境のログ表示
logs-prod:
	docker-compose logs -f

# コンテナの再起動
restart:
	docker-compose restart

# 開発環境のコンテナに入る
shell:
	docker-compose -f docker-compose.dev.yml exec api sh

# 本番環境のコンテナに入る
shell-prod:
	docker-compose exec api sh

# データベースコンテナに入る
db-shell:
	docker-compose -f docker-compose.dev.yml exec postgres psql -U postgres -d stackies_dev

# データベースをリセット
db-reset:
	docker-compose -f docker-compose.dev.yml down -v
	docker-compose -f docker-compose.dev.yml up -d postgres 
