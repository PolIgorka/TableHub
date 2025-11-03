# build stage
FROM golang:1.23-alpine AS builder
WORKDIR /app

# устанавливаем git и сертификаты
RUN apk update && apk add --no-cache git ca-certificates

# копируем модули и скачиваем зависимости
COPY go.mod go.sum ./
RUN go mod download

# копируем исходники
COPY . .

# собираем статический бинарник
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /bin/tablehub-server ./cmd/server

# runtime stage
FROM scratch

# копируем сертификаты (для HTTPS)
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
COPY --from=builder /bin/tablehub-server /bin/tablehub-server
COPY config.yaml /config.yaml

EXPOSE 8081

ENTRYPOINT ["/bin/tablehub-server"]
