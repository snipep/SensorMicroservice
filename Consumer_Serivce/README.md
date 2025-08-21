## Consumer Service

REST + gRPC service that stores and serves sensor data, and provides user auth (signup/signin) with JWT. Incoming gRPC data is queued on a buffered channel and inserted by a worker pool.

### Features

- User signup/signin with bcrypt password hashing
- JWT issuance (HS256) and global auth middleware (all routes except `/api/v1/signup` and `/api/v1/signin`)
- Sensor data query/edit/delete APIs over HTTP
- gRPC server for ingesting sensor streams
- Buffered channel + worker pool for DB inserts (configurable)

### Requirements

- Go 1.24+
- MySQL 8+

### Environment variables

- `ADDR` (default `:8080`): HTTP listen address
- `GRPC_PORT` (default `:50051`): gRPC listen address
- `JWT_SECRET` (default `dev-secret-change`): HMAC secret for JWT
- `DB_ADDR` (default `user:password@tcp(127.0.0.1:3306)/sensordata?charset=utf8mb4&parseTime=true&loc=UTC`): MySQL DSN
- `INGEST_WORKERS` (default `4`): number of goroutine workers for inserts
- `INGEST_BUFFER` (default `1024`): buffered channel capacity

Recommended (no quotes):

```bash
export DB_ADDR=user:password@tcp(127.0.0.1:3306)/sensordata?charset=utf8mb4&parseTime=true&loc=UTC
export ADDR=:8080 GRPC_PORT=:50051 JWT_SECRET=change-me INGEST_WORKERS=8 INGEST_BUFFER=4096
```

### Database & migrations

- SQL lives in `cmd/migrate/migrations`
- On service startup, embedded `.up.sql` files run automatically
- `users.created_at` uses `DATETIME`; the server writes UTC timestamps

If you use `golang-migrate` manually:

```bash
migrate -path cmd/migrate/migrations -database "$DB_URL" up
```

### Run locally

```bash
go build -o ./bin/consumer ./cmd/api
ADDR=:8080 GRPC_PORT=:50051 JWT_SECRET=change-me DB_ADDR="user:password@tcp(127.0.0.1:3306)/sensordata?charset=utf8mb4&parseTime=true&loc=UTC" \
  INGEST_WORKERS=4 INGEST_BUFFER=1024 \
  ./bin/consumer
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
  -e INGEST_WORKERS=8 -e INGEST_BUFFER=4096 \
  --name consumer consumer-service:latest
```

### Docker Compose

From repo root:

```bash
make up         # or: docker compose up --build -d
make logs       # tail logs
```

### API Documentation (Swagger) 
This service uses Swagger to provide interactive API documentation. You can use this interface to view all available endpoints, see their parameters, and execute API requests directly from your browser.

Once the server is running, access the Swagger UI at:

`http://localhost:8080/swagger/index.html`

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

### Worker pool

- gRPC handler enqueues sensor points to a buffered channel of size `INGEST_BUFFER`
- `INGEST_WORKERS` goroutines consume and insert into MySQL
- Provides burst handling and backpressure while maximizing throughput
