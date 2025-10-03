# Gunakan versi Go yang sesuai requirement
FROM golang:1.24 AS builder

WORKDIR /app

# Copy go.mod dan go.sum dulu
COPY go.mod go.sum ./
RUN go mod download

# Copy semua source code
COPY . .

# Build binary dari folder yang ada main.go
RUN go build -o main ./cmd/server

# Stage kedua, image lebih kecil
FROM debian:bookworm-slim
WORKDIR /app
COPY --from=builder /app/main .

CMD ["./main"]
