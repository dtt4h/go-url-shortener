# URL Shortener - Руководство по развёртыванию

## Содержание
- [Быстрый старт](#быстрый-старт)
- [Архитектура](#архитектура)
- [Проверка работоспособности](#проверка-работоспособности)
- [API Endpoints](#api-endpoints)
- [Устранение проблем](#устранение-проблем)

---

## Быстрый старт

### Предварительные требования
- Docker
- Docker Compose

### Запуск

```bash
# Клонировать репозиторий и перейти в директорию проекта
cd go-url-shortener

# Запустить все сервисы
docker compose up -d --build

# Проверить статус контейнеров
docker compose ps
```

### Остановка

```bash
# Остановить все сервисы
docker compose down

# Остановить и удалить volumes
docker compose down -v
```

---

## Архитектура

Проект состоит из 5 сервисов:

| Сервис | Порт | Описание |
|--------|------|----------|
| PostgreSQL | 5432 | База данных |
| Zookeeper | 2181 | Координация Kafka |
| Kafka | 9092 | Очередь сообщений |
| Backend (Go) | 8080 | API сервер |
| Frontend (Nginx) | 80 | Веб-интерфейс + прокси |

### Структура проксирования Nginx

```
localhost:80/           -> Frontend (index.html)
localhost:80/api/*      -> Backend (проксирование)
localhost:80/:code      -> Backend (редирект на оригинальный URL)
```

---

## Проверка работоспособности

### 1. Проверка фронтенда

```bash
# Должен вернуть HTML страницу
curl http://localhost:80/
```

Ожидаемый вывод: HTML с формой сокращения URL.

### 2. Проверка API

```bash
# Создание короткой ссылки
curl -X POST http://localhost:80/api/v1/shorten \
  -H "Content-Type: application/json" \
  -d '{"url":"https://google.com"}'
```

Ожидаемый ответ:
```json
{"short_url":"http://localhost:8080/abc123de","expires_at":0,"click_count":0}
```

### 3. Проверка редиректа

```bash
# Проверка редиректа (возвращает 301)
curl -I http://localhost:80/abc123de
```

### 4. Проверка QR-кода

```bash
# Проверка генерации QR-кода (возвращает изображение)
curl -I http://localhost:80/abc123de/qr
```

### 5. Проверка через браузер

1. Открыть `http://localhost:80` в браузере
2. Ввести URL в поле ввода
3. Нажать "Shorten"
4. Скопировать полученную ссылку или QR-код

---

## API Endpoints

| Метод | Путь | Описание |
|-------|------|----------|
| POST | `/api/v1/shorten` | Создать короткую ссылку |
| GET | `/:code` | Перейти по короткой ссылке |
| DELETE | `/:code` | Удалить ссылку |
| GET | `/:code/qr` | Получить QR-код |

### POST /api/v1/shorten

**Request:**
```json
{
  "url": "https://example.com/very/long/url"
}
```

**Response:**
```json
{
  "short_url": "http://localhost:8080/abc123de",
  "expires_at": 0,
  "click_count": 0
}
```

---

## Устранение проблем

### Контейнеры не запускаются

```bash
# Просмотр логов
docker compose logs

# Перезапуск
docker compose restart
```

### Ошибка миграций

```bash
# Применить миграции вручную
docker exec url-shortener-db psql -U postgres -d url_shortener -f /docker-entrypoint-initdb.d/001_create_urls.sql
```

### Очистка и пересборка

```bash
# Полная очистка
docker compose down -v
docker system prune -f

# Пересборка
docker compose up -d --build
```

### Проверка БД

```bash
# Подключение к PostgreSQL
docker exec -it url-shortener-db psql -U postgres -d url_shortener

# Просмотр таблицы ссылок
SELECT * FROM urls;
```

### Проверка Kafka

```bash
# Подключение к Kafka
docker exec -it url-shortener-kafka kafka-console-consumer --bootstrap-server localhost:9092 --topic url-events --from-beginning
```

---

## Переменные окружения

| Переменная | Описание | Значение по умолчанию |
|------------|----------|----------------------|
| `CONFIG_PATH` | Путь к конфигу | `configs/config_local.yaml` |
| `DB_URL` | Строка подключения к БД | `postgres://postgres:postgres@postgres:5432/url_shortener` |
| `KAFKA_BROKERS` | Kafka брокеры | `kafka:9092` |
| `URL_BASE` | Базовый URL | `http://localhost:8080/` |

---

## Команды Makefile

```bash
make help     # Показать все команды
make up       # Запустить (фоновый режим)
make down     # Остановить
make logs     # Логи всех сервисов
make test     # Запустить тесты
make clean    # Очистить контейнеры и volumes
```
