# URL Shortener

Сервис для сокращения URL-адресов.

## Технологии

- **Go** — язык программирования
- **PostgreSQL** — база данных
- **Gin** — веб-фреймворк
- **Kafka** — брокер сообщений

## Требования

- Go 1.21+
- PostgreSQL 14+
- Kafka 3.0+

## Установка

```bash
# Клонирование репозитория
git clone <repository-url>
cd <project-name>

# Установка зависимостей
go mod download
```

## Конфигурация

Настройки хранятся в файле `.env`:

```env
# База данных
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=password
DB_NAME=url_shortener

# Kafka
KAFKA_BROKERS=localhost:9092

# Приложение
APP_HOST=0.0.0.0
APP_PORT=8080
```

## Запуск

```bash
# Запуск БД и Kafka через Docker
docker-compose up -d

# Запуск приложения
go run cmd/main.go
```

## API

| Метод | Эндпоинт | Описание |
|-------|----------|----------|
| POST | `/api/urls` | Создать короткую ссылку |
| GET | `/api/urls/:short_code` | Получить оригинальный URL |
| GET | `/:short_code` | Перенаправить на оригинальный URL |
| DELETE | `/api/urls/:short_code` | Удалить ссылку |

### Примеры запросов

```bash
# Создание ссылки
curl -X POST http://localhost:8080/api/urls \
  -H "Content-Type: application/json" \
  -d '{"original_url": "https://example.com"}'

# Ответ
{"short_code":"abc123","original_url":"https://example.com","created_at":"2024-01-01T00:00:00Z"}
```

## Структура проекта

```
├── cmd/
│   └── main.go           # Точка входа
├── internal/
│   ├── config/           # Конфигурация
│   ├── handler/          # HTTP-обработчики
│   ├── model/            # Модели данных
│   ├── repository/       # Работа с БД
│   └── service/          # Бизнес-логика
├── migrations/           # Миграции БД
├── docker-compose.yml    # Docker-конфигурация
├── .env                  # Переменные окружения
└── go.mod                # Зависимости
```

## Тесты

```bash
go test ./...
```

## Лицензия

MIT
