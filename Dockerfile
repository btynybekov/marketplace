# Stage 1: Build
FROM golang:1.21-alpine AS builder

WORKDIR /app

# Копируем go.mod и go.sum
COPY go.mod go.sum ./
RUN go mod download

# Копируем исходники
COPY . .

# Собираем бинарник
RUN go build -o marketplace ./cmd/server

# Stage 2: Run
FROM alpine:latest

WORKDIR /app

# Копируем бинарник из build stage
COPY --from=builder /app/marketplace .

# Копируем configs (опционально)
COPY configs ./configs

# Устанавливаем timezone (по желанию)
RUN apk add --no-cache tzdata

# Запуск приложения
CMD ["./marketplace"]