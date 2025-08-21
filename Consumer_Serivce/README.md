## Consumer Service

REST + gRPC service that stores and serves sensor data, and provides user auth (signup/signin) with JWT protection on all routes except auth.

### Features

- User signup/signin with bcrypt password hashing
- JWT issuance (HS256) and global auth middleware (all routes except `/api/v1/signup` and `/api/v1/signin`)
- Sensor CRUD-ish APIs over HTTP
- gRPC server for ingesting sensor streams

### Requirements

- Go 1.24+
- MySQL 8+

### Environment variables

- `ADDR` (default `:8080`): HTTP listen address
- `GRPC_PORT` (default `:50051`): gRPC listen address
- `JWT_SECRET` (default `dev-secret-change`): HMAC secret for JWT
- `DB_ADDR` (default `user:password@tcp(127.0.0.1:3306)/sensordata?charset=utf8mb4&parseTime=true&loc=UTC`): MySQL DSN

Recommended (no quotes):

```bash
export DB_ADDR=user:password@tcp(127.0.0.1:3306)/sensordata?charset=utf8mb4&parseTime=true&loc=UTC
export ADDR=:8080 GRPC_PORT=:50051 JWT_SECRET=change-me
```

### Database & migrations

Schema lives in `cmd/migrate/migrations`. The `users` table uses `DATETIME` and the service writes `created_at` in UTC.

If you use `golang-migrate`:

```bash
migrate -path cmd/migrate/migrations -database "$DB_URL" up
```

### Run locally

```bash
go build -o ./bin/consumer ./cmd/api
ADDR=:8080 GRPC_PORT=:50051 JWT_SECRET=change-me DB_ADDR="user:password@tcp(127.0.0.1:3306)/sensordata?charset=utf8mb4&parseTime=true&loc=UTC" ./bin/consumer
```

### Docker

Build:

```bash
docker build -t consumer-service:latest -f Dockerfile .
```

Run:

```bash
docker run --rm -p 8080:8080 -p 50051:50051 \
  -e ADDR=:8080 -e GRPC_PORT=:50051 \
  -e JWT_SECRET=change-me \
  -e DB_ADDR='user:password@tcp(db:3306)/sensordata?charset=utf8mb4&parseTime=true&loc=UTC' \
  --name consumer consumer-service:latest
```

### Docker Compose

From repo root:

```bash
docker compose up --build
```

### HTTP API

Base path: `/api/v1`

Auth (public):

```bash
# Signup
curl -X POST http://localhost:8080/api/v1/signup \
  -H 'Content-Type: application/json' \
  -d '{"name":"Alice","email":"alice@example.com","password":"P@ssw0rd123"}'

# Signin -> { token }
curl -X POST http://localhost:8080/api/v1/signin \
  -H 'Content-Type: application/json' \
  -d '{"email":"alice@example.com","password":"P@ssw0rd123"}'
```

Protected (JWT required in `Authorization: Bearer <token>`):

```bash
curl -X PUT http://localhost:8080/api/v1/sensordata/query \
  -H 'Authorization: Bearer <TOKEN>' -H 'Content-Type: application/json' \
  -d '{"id1":"A","id2":1001,"new_value":42.0}'
```

Other endpoints (examples):

- `GET /api/v1/sensordata/query`
- `GET /api/v1/sensordata/history`
- `GET /api/v1/sensordata/query-history`
- `DELETE /api/v1/sensordata/...`
- `PUT /api/v1/sensordata/...` (JWT)

### gRPC

Listens on `GRPC_PORT` (default `:50051`). See `Proto/` for service definitions.
