# Quick start

## Requirements
- Docker Desktop / Docker Engine
- Ports `8080` и `5432` are free

## Run locally (docker-compose)
```bash
docker compose up -d --build
docker compose logs -f server
curl http://localhost:8080/healthz   # -> ok
```

## Manual server start
```bash
export DATABASE_URL="postgres://anti:anti@localhost:5432/anti?sslmode=disable"
cd server
go run ./cmd/server
```

## Enabling the Mailer Node (HTTP API)
In environment variables:
```
NODE_ENABLED=true
REQUIRE_ACK_SIGNATURE=true   # optional (recommended)
```
Endpoints will appear under `/node/*`. See "Node Mailbox".