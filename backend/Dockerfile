FROM golang:1.25-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o /app/main ./cmd/main.go

FROM alpine:3.19

WORKDIR /app

COPY --from=builder /app/main .
COPY configs/ ./configs/

EXPOSE 8080

CMD ["./main"]
