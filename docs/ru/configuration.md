# Конфигурация

## Переменные окружения (сервер)

| Переменная | По умолчанию | Описание |
|---|---:|---|
| `DATABASE_URL` | — | DSN Postgres (обязательно) |
| `ADDR` | `:8080` | Адрес HTTP/WS |
| `PENDING_BATCH` | `1000` | Сколько pending отдаётся за раз |
| `RETAIN_DELIVERED_MINUTES` | `60` | TTL доставленных |
| `RETAIN_UNDELIVERED_DAYS` | `14` | TTL недоставленных |
| `ALLOWED_ORIGINS` | пусто | Разрешённые Origin (через запятую) |
| `JANITOR_INTERVAL_SECONDS` | `60` | Период фоновой очистки |
| `READ_TIMEOUT_SECONDS` | `60` | WS read timeout |
| `WRITE_TIMEOUT_SECONDS` | `20` | WS write timeout |
| `HANDSHAKE_TIMEOUT_SECONDS` | `5` | HTTP read-header timeout |
| `NODE_ENABLED` | `false` | Включить HTTP-узел почтового ящика |
| `NODE_PEERS` | пусто | Соседние узлы (через запятую) для репликации |
| `REQUIRE_ACK_SIGNATURE` | `false` | Требовать подпись ACK (Ed25519) |

## Пример `.env.example`
```env
DATABASE_URL=postgres://anti:anti@db:5432/anti?sslmode=disable
ADDR=:8080
PENDING_BATCH=1000
RETAIN_DELIVERED_MINUTES=60
RETAIN_UNDELIVERED_DAYS=14
ALLOWED_ORIGINS=
JANITOR_INTERVAL_SECONDS=60
READ_TIMEOUT_SECONDS=60
WRITE_TIMEOUT_SECONDS=20
HANDSHAKE_TIMEOUT_SECONDS=5
NODE_ENABLED=true
NODE_PEERS=
REQUIRE_ACK_SIGNATURE=true
```
