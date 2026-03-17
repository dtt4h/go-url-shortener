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

### Эндпоинты

| Метод   | Эндпоинт              | Описание                              |
|---------|-----------------------|---------------------------------------|
| POST    | `/api/v1/shorten`     | Создать короткую ссылку               |
| GET     | `/:code`              | Перенаправить на оригинальный URL     |
| GET     | `/:code/qr`           | Получить QR-код для короткой ссылки   |
| DELETE  | `/:code`              | Удалить ссылку                        |
| GET     | `/`                   | Главная страница (HTML)               |

---

### POST /api/v1/shorten

Создать короткую ссылку.

**Request:**

```json
{
  "url": "https://example.com/very/long/path"
}
```

**Response (201 Created):**

```json
{
  "short_url": "http://localhost:8080/abc123",
  "expires_at": 0,
  "click_count": 0
}
```

**Ошибки:**

- `400 Bad Request` — некорректный URL
- `500 Internal Server Error` — внутренняя ошибка сервера

---

### GET /:code

Перенаправить на оригинальный URL.

**Параметры пути:**

- `code` — короткий код ссылки

**Response:**

- `301 Moved Permanently` — редирект на оригинальный URL

**Ошибки:**

- `404 Not Found` — ссылка не найдена
- `410 Gone` — ссылка истекла

---

### GET /:code/qr

Получить QR-код для короткой ссылки (PNG-изображение).

**Параметры пути:**

- `code` — короткий код ссылки

**Response:**

- `200 OK` — изображение PNG (256x256)

**Ошибки:**

- `404 Not Found` — ссылка не найдена
- `410 Gone` — ссылка истекла
- `500 Internal Server Error` — ошибка генерации QR-кода

**Пример использования:**

```html
<img src="http://localhost:8080/abc123/qr" alt="QR код" />
```

---

### DELETE /:code

Удалить ссылку.

**Параметры пути:**

- `code` — короткий код ссылки

**Response:**

- `204 No Content` — ссылка успешно удалена

**Ошибки:**

- `404 Not Found` — ссылка не найдена

---

### GET /

Главная страница.

**Response:**

- `200 OK` — HTML-страница

---

## Примеры (Postman)

### Создание ссылки

- **Method:** POST
- **URL:** `http://localhost:8080/api/v1/shorten`
- **Headers:** `Content-Type: application/json`
- **Body (JSON):**
```json
{
  "url": "https://example.com/very/long/path"
}
```

**Response (201):**
```json
{
  "short_url": "http://localhost:8080/abc123",
  "expires_at": 0,
  "click_count": 0
}
```

---

### Переход по короткой ссылке

- **Method:** GET
- **URL:** `http://localhost:8080/abc123`

**Response:** 301 Moved Permanently → редирект на оригинальный URL

---

### Получение QR-кода

- **Method:** GET
- **URL:** `http://localhost:8080/abc123/qr`

**Response:** 200 OK (image/png)

---

### Удаление ссылки

- **Method:** DELETE
- **URL:** `http://localhost:8080/abc123`

**Response:** 204 No Content

## Структура проекта

```
├── cmd/
│   └── main.go           # Точка входа
├── internal/
│   ├── config/           # Конфигурация
│   ├── handler/          # HTTP-обработчики
│   ├── middleware/       # Промежуточное ПО
│   ├── model/            # Модели данных
│   ├── repository/       # Работа с БД
│   └── service/          # Бизнес-логика
├── migrations/           # Миграции БД
├── web/                  # Статические файлы
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