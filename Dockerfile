FROM golang:1.23.6 AS builder

# Установите рабочую директорию
WORKDIR /app

# Копируйте go.mod и go.sum и загрузите зависимости
COPY go.mod go.sum ./
RUN go mod download

# Копируйте исходный код приложения
COPY . .

# Соберите приложение
RUN go build -o main .

# Используйте легкий образ для выполнения приложения
FROM alpine:latest

# Установите необходимые зависимости для выполнения приложения
RUN apk --no-cache add ca-certificates

# Копируйте исполняемый файл из предыдущего этапа
COPY --from=builder /app/main .

# Убедитесь, что приложение слушает на порту 8080
EXPOSE 8080

# Запустите приложение
CMD ["./main"]