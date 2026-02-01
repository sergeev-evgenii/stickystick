# ------------------------------
# Stage 1: Build
# ------------------------------
FROM golang:1.24.0 AS builder

WORKDIR /app

# Копируем только go.mod и go.sum для кеширования зависимостей
COPY go.mod go.sum ./

# Загружаем зависимости
RUN go mod download

# Копируем весь исходный код, кроме папки uploads
COPY . .

# Собираем статический бинарь
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o backend ./cmd/api

# ------------------------------
# Stage 2: Runtime
# ------------------------------
FROM gcr.io/distroless/base-debian11

WORKDIR /app

# Копируем только бинарь
COPY --from=builder /app/backend .

# Экспонируем порт
EXPOSE 5000

# Команда запуска
CMD ["./backend"]