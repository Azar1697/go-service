# --- Stage 1: Сборка (Builder) ---
FROM golang:1.23-alpine AS builder

WORKDIR /app

# Кэшируем зависимости (чтобы не качать их каждый раз)
COPY go.mod go.sum ./
RUN go mod download

# Копируем исходный код и собираем
COPY . .
RUN go build -o my-service .

# --- Stage 2: Финальный образ (Runner) ---
FROM alpine:latest

WORKDIR /app

# Копируем только скомпилированный файл из первого этапа
COPY --from=builder /app/my-service .

# Открываем порт
EXPOSE 8080

# Запускаем
CMD ["./my-service"]
