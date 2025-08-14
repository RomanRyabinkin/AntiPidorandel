# Configuration

## Environment variables (server)

| Variable | Default | Description |
|---|---:|---|
| `DATABASE_URL` | — | DSN Postgres (necessarily) |
| `ADDR` | `:8080` | HTTP/WS address |
| `PENDING_BATCH` | `1000` | How much pending is given at a time |
| `RETAIN_DELIVERED_MINUTES` | `60` | TTL delivered |
| `RETAIN_UNDELIVERED_DAYS` | `14` | TTL undelivered |
| `ALLOWED_ORIGINS` | empty | Allowed Origins (comma separated) |
| `JANITOR_INTERVAL_SECONDS` | `60` | Background cleaning period |
| `READ_TIMEOUT_SECONDS` | `60` | WS read timeout |
| `WRITE_TIMEOUT_SECONDS` | `20` | WS write timeout |
| `HANDSHAKE_TIMEOUT_SECONDS` | `5` | HTTP read-header timeout |
| `NODE_ENABLED` | `false` | Enable HTTP Mailbox Host |
| `NODE_PEERS` | empty | Neighboring nodes (comma separated) for replication |
| `REQUIRE_ACK_SIGNATURE` | `false` | Require ACK signature (Ed25519) |

## Example `.env.example`
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