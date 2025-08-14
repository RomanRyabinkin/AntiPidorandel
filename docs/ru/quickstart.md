# Быстрый старт

## Требования
- Docker Desktop / Docker Engine
- Порты `8080` и `5432` свободны

## Запуск локально (docker-compose)
```bash
docker compose up -d --build
docker compose logs -f server
curl http://localhost:8080/healthz   # -> ok
```

## Ручной запуск сервера
```bash
export DATABASE_URL="postgres://anti:anti@localhost:5432/anti?sslmode=disable"
cd server
go run ./cmd/server
```

## Включение узла-почтовика (HTTP API)
В переменных окружения:
```
NODE_ENABLED=true
REQUIRE_ACK_SIGNATURE=true   # опционально (рекомендуется)
```
Эндпоинты появятся под `/node/*`. См. «Узлы‑почтовики».
