FROM golang:1.27.1 AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . .
# Собираем приложение с включенным CGO (обязательно для SQLite)
RUN CGO_ENABLED=0 GOOS=linux go build -o scheduler main.go

FROM alpine:latest

WORKDIR /app
# Копируем бинарник и статические файлы из сборочного образа
COPY --from=builder /app/scheduler .
COPY web ./web

# Указываем дефолтные переменные окружения
ENV TODO_PORT=7540
ENV TODO_DBFILE=scheduler.db
ENV TODO_PASSWORD=""

EXPOSE ${TODO_PORT}

CMD ["./scheduler"]