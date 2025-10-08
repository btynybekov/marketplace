# Stage 1: builder
FROM golang:1.25-alpine AS builder

WORKDIR /app
RUN apk add --no-cache git bash
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o marketplace

# Stage 2: final
FROM alpine:3.18
WORKDIR /app

# Копируем бинарь и шаблоны
COPY --from=builder /app/marketplace .
COPY --from=builder /app/templates ./templates

# Копируем .env внутрь контейнера
COPY --from=builder /app/.env ./

# Устанавливаем зависимости
RUN apk add --no-cache bash ca-certificates

# Запуск приложения
CMD ["./marketplace"]