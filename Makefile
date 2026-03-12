.PHONY: help up up-d down restart build ps logs logs-backend logs-frontend logs-db clean prune \
	migrate migrate-create db-connect backend-shell frontend-shell test lint start

help:
	@echo "URL Shortener - Команды:"
	@echo ""
	@echo "  Docker:"
	@echo "    make start        - поднять БД, миграции, приложение"
	@echo "    make up           - поднять (фоновый режим)"
	@echo "    make up-run       - поднять (интерактивный режим)"
	@echo "    make down         - остановить"
	@echo "    make restart      - перезапустить"
	@echo "    make build        - пересобрать"
	@echo "    make ps           - статус контейнеров"
	@echo ""
	@echo "  Логи:"
	@echo "    make logs         - все сервисы"
	@echo "    make logs-backend - backend"
	@echo "    make logs-frontend - frontend"
	@echo "    make logs-db      - PostgreSQL"
	@echo ""
	@echo "  База данных:"
	@echo "    make migrate      - применить миграции"
	@echo "    make migrate-create NAME=... - создать миграцию"
	@echo "    make db-connect   - подключиться к БД"
	@echo ""
	@echo "  Шеллы:"
	@echo "    make backend-shell  - зайти в backend"
	@echo "    make frontend-shell - зайти в frontend"
	@echo ""
	@echo "  Тесты/Линтеры:"
	@echo "    make test        - тесты"
	@echo "    make lint        - линтер"
	@echo ""
	@echo "  Очистка:"
	@echo "    make clean       - удалить контейнеры и volumes"
	@echo "    make prune       - полная очистка"

# Docker
up:
	docker compose -f docker-compose.yml up -d

up-run:
	docker compose -f docker-compose.yml up

down:
	docker compose -f docker-compose.yml down

restart:
	docker compose -f docker-compose.yml down && docker compose -f docker-compose.yml up -d

build:
	docker compose -f docker-compose.yml build --no-cache --force-rm

ps:
	docker compose -f docker-compose.yml ps

# Логи
logs:
	docker compose -f docker-compose.yml logs -f

logs-backend:
	docker compose -f docker-compose.yml logs -f backend

logs-frontend:
	docker compose -f docker-compose.yml logs -f frontend

logs-db:
	docker compose -f docker-compose.yml logs -f postgres

# Миграции
migrate-up:
	docker cp backend/migrations/001_create_urls.sql url-shortener-db:/tmp/
	docker exec -it url-shortener-db psql -U postgres -d url_shortener -f /tmp/001_create_urls.sql && \
	docker exec url-shortener-db rm /tmp/001_create_urls.sql

migrate-down:
	docker exec -it url-shortener-db psql -U postgres -d url_shortener -c "DROP TABLE IF EXISTS urls CASCADE;"

migrate: migrate-up

start:
	docker compose up -d --build


migrate-create:
ifndef NAME
	@echo "Ошибка: укажите имя: make migrate-create NAME=create_users"
	@exit 1
endif
	docker run --rm -it \
		-v $(PWD)/backend/migrations:/migrations \
		migrate/migrate:latest \
		create -dir /migrations -seq -ext sql $(NAME)

# База данных
db-connect:
	docker exec -it url-shortener-db psql -U postgres -d url_shortener

# Шеллы
backend-shell:
	docker exec -it url-shortener-backend sh

frontend-shell:
	docker exec -it url-shortener-frontend sh

# Тесты и линтеры
test:
	docker exec url-shortener-backend go test -v ./...

lint:
	docker exec url-shortener-backend golangci-lint run ./...

# Очистка
clean:
	docker compose -f docker-compose.yml down -v

prune:
	docker compose -f docker-compose.yml down -v --rmi local
	docker system prune -f