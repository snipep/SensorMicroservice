# syntax=docker/dockerfile:1

FROM golang:1.24-alpine AS builder
WORKDIR /src
RUN apk add --no-cache ca-certificates tzdata git

# Cache modules
COPY go.mod go.sum ./
RUN go mod download

# Copy source
COPY . .

# Build static binary
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /out/consumer-service ./cmd/api

FROM alpine:3.20
WORKDIR /app
RUN apk add --no-cache ca-certificates tzdata

COPY --from=builder /out/consumer-service /app/consumer-service

# Default envs (override at runtime)
ENV ADDR=:8080 \
    GRPC_PORT=:50051 \
    JWT_SECRET=dev-secret-change \
    DB_ADDR="user:password@tcp(db:3306)/sensordata?charset=utf8mb4&parseTime=true&loc=UTC"

EXPOSE 8080 50051

CMD ["/app/consumer-service"]


