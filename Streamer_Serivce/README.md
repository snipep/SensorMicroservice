## Streamer Service

Generates and streams sensor data to the Consumer service via gRPC, and exposes an HTTP endpoint to adjust streaming frequency. Multiple streamer instances can run concurrently, each with a fixed sensor ID/type.

### Features

- gRPC client to Consumer (`MICROSERVICE_B_ADDR`)
- Fixed sensor configuration per instance via env: `STREAMER_ID1`, `STREAMER_TYPE`
- HTTP endpoint to update frequency: `POST /api/v1/frequency` with `{ "frequency_seconds": 5 }`

### Requirements

- Go 1.24+

### Environment variables

- `MICROSERVICE_B_ADDR` (default `localhost:50051`): Consumer gRPC address
- `HTTP_PORT` (default `:8080`): HTTP listen port
- `STREAMER_ID1` (default `A`): sensor key
- `STREAMER_TYPE` (default `Temperature`): sensor type label

### Run locally

```bash
go build -o ./bin/streamer ./cmd
MICROSERVICE_B_ADDR=localhost:50051 HTTP_PORT=:8081 STREAMER_ID1=B STREAMER_TYPE=Humidity ./bin/streamer
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
  -e STREAMER_ID1=B -e STREAMER_TYPE=Humidity \
  --name streamer-b streamer-service:latest
```

### Docker Compose

The repo-level compose runs three streamers by default:

- `temperature_streamer` (A/Temperature) on 8081
- `humidity_streamer` (B/Humidity) on 8082
- `pressure_streamer` (C/Pressure) on 8083

From repo root:

```bash
make up      # or: docker compose up --build -d
```

### HTTP API

Adjust frequency (per instance):

```bash
curl -X POST http://localhost:8081/api/v1/frequency \
  -H 'Content-Type: application/json' \
  -d '{"frequency_seconds": 5}'
```
