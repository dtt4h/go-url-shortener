FROM golang:1.25-alpine AS builder

WORKDIR /app

COPY backend/go.mod backend/go.sum ./
RUN go mod download

COPY backend/ .

RUN CGO_ENABLED=0 GOOS=linux go build -o /app/main ./cmd/main.go

FROM alpine:3.19

WORKDIR /app

RUN apk add --no-cache ca-certificates

# Создание не-root пользователя
RUN addgroup -g 1000 appgroup && \
    adduser -u 1000 -G appgroup -s /bin/sh -D appuser

# Создание директорий с правильными правами
RUN mkdir -p /app/configs /app/web && \
    chown -R appuser:appgroup /app

USER appuser

COPY --from=builder --chown=appuser:appgroup /app/main .
COPY --chown=appuser:appgroup backend/configs/ ./configs/
COPY --chown=appuser:appgroup frontend/web/ ./web/

EXPOSE 8080

HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD wget --no-verbose --tries=1 --spider http://localhost:8080/health || exit 1

CMD ["./main"]
