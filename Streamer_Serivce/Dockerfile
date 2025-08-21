# syntax=docker/dockerfile:1

FROM golang:1.24-alpine AS builder
WORKDIR /src
RUN apk add --no-cache ca-certificates tzdata git

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /out/streamer-service ./cmd

FROM alpine:3.20
WORKDIR /app
RUN apk add --no-cache ca-certificates tzdata

COPY --from=builder /out/streamer-service /app/streamer-service

# Defaults (override in runtime)
ENV MICROSERVICE_B_ADDR=localhost:50051 \
    HTTP_PORT=:8081

EXPOSE 8081

CMD ["/app/streamer-service"]


