# Builder
FROM golang:1.24-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o icu-app cmd/main/app.go
RUN go build -o migrate-tool cmd/migrate/main.go

# Runnable
FROM alpine:latest
WORKDIR /root/
COPY --from=builder /app/icu-app .
COPY --from=builder /app/migrate-tool .
# Copy migration files for the tool to use
COPY --from=builder /app/etc/migrations ./etc/migrations

CMD ["./icu-app"]
