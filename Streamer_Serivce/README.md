## Streamer Service

Generates and streams sensor data to the Consumer service via gRPC, and exposes an HTTP endpoint to adjust streaming frequency.

### Features

- gRPC client to `consumer` service to send sensor data
- HTTP endpoint to update frequency: `POST /api/v1/frequency` with `{ "frequency_seconds": 5 }`

### Requirements

- Go 1.24+

### Environment variables

- `MICROSERVICE_B_ADDR` (default `localhost:50051`): Consumer gRPC address
- `HTTP_PORT` (default `:8080`): HTTP listen port for control endpoint

### Run locally

```bash
go build -o ./bin/streamer ./cmd
MICROSERVICE_B_ADDR=localhost:50051 HTTP_PORT=:8081 ./bin/streamer
```

### Docker

Build:

```bash
docker build -t streamer-service:latest -f Dockerfile .
```

Run:

```bash
docker run --rm -p 8081:8080 \
  -e MICROSERVICE_B_ADDR=host.docker.internal:50051 \
  -e HTTP_PORT=:8080 \
  --name streamer streamer-service:latest
```

### Docker Compose

From repo root:

```bash
docker compose up --build
```

In compose, the service uses `MICROSERVICE_B_ADDR=consumer:50051` to reach the consumer container.

### HTTP API

```bash
curl -X POST http://localhost:8081/api/v1/frequency \
  -H 'Content-Type: application/json' \
  -d '{"frequency_seconds": 5}'
```
