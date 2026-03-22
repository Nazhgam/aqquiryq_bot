FROM golang:1.24-alpine AS builder

WORKDIR /app

# Установка необходимых системных зависимостей (если понадобятся для сборки)
RUN apk add --no-cache git

COPY go.mod go.sum ./
RUN go mod download

COPY . .

# Сборка приложения
RUN go build -o app ./cmd/server/main.go

FROM alpine:latest

WORKDIR /app

# Установка CA-сертификатов для HTTPS-запросов (важно для Telegram API)
RUN apk add --no-cache ca-certificates

COPY --from=builder /app/app .
COPY --from=builder /app/config ./config
COPY --from=builder /app/templates ./templates

CMD ["./app"]
