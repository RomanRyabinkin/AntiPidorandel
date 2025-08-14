# AntiPidorandel: Private E2EE Messenger

## 🇬🇧 English

> **Default branch**: `master` (all actual code and documentation live in this branch).

**AntiPidorandel** is a privacy-first, invite-only messenger. The server acts as a **blind relay**: it stores and forwards only **encrypted envelopes** (E2EE). Content is never decrypted server-side, and server metadata is minimized.

### Documentation

Link: https://romanryabinkin.github.io/AntiPidorandel/en/


---



### Features

* **Blind Relay Server (Go)**
  * WebSocket API: `WS /ws?user_id=<id>`
  * Immediate delivery to online peers, queued delivery for offline peers
  * Delivery acknowledgments (`ack`)
* **Encrypted Storage (PostgreSQL)**
  * Stores only ciphertext + minimal routing metadata (`to_user`, timestamps, TTL)
  * Automatic cleanup: delivered and expired messages
* **Ops & Runtime**
  * Health probe: `GET /healthz`
  * **Docker** / **docker-compose** (distroless runtime)
  * Env-based configuration
* **Security Posture**
  * E2EE by design (server never sees plaintext)
  * Minimal logs/metadata

---

### Project Structure

* **server/** — standalone Go module (the server)
  * `cmd/server/` — entrypoint
  * `internal/transport/ws/` — WebSocket transport
  * `internal/store/` — storage interface
    * `postgres/` — pgx implementation
  * `internal/hub/` — online session hub
  * `internal/config/` — env config
  * `internal/wire/` — on-the-wire frame contracts (JSON)
  * `DockerFile` — multi-stage build (distroless runtime)
  * `go.mod`, `go.sum`, `.gitignore`, `.dockerignore`
* **docker-compose.yml** (repo root) — local Postgres + server
* **go.work** (repo root) — Go workspace, includes `./server`

---

### Technologies & Libraries

* Go 1.23+ (`gorilla/websocket`, `pgx`)
* PostgreSQL 16+
* Docker / docker-compose
* Distroless runtime (non-root)
* Recommended IDE: VS Code (Go, Docker, YAML, SQLTools, Error Lens)

---

### Quick Start

**Requirements:**
* Docker Desktop / Docker Engine
* Free ports `8080` (server) and `5432` (Postgres)

**Run with docker-compose (recommended):**
```bash
docker compose up -d --build
docker compose logs -f server
curl http://localhost:8080/healthz   # -> ok
```

**Manual run (no Docker):**
```bash
# Start Postgres and create DB "anti"
export DATABASE_URL="postgres://anti:anti@localhost:5432/anti?sslmode=disable"
cd server
go run ./cmd/server
```

---

### Configuration

**Environment variables (server):**

* **DATABASE_URL** — Postgres DSN (required)

* **ADDR** — HTTP/WS listen address (default :8080)

* **PENDING_BATCH** — pending messages per fetch (default 1000)

* **RETAIN_DELIVERED_MINUTES** — TTL for delivered (default 60)

* **RETAIN_UNDELIVERED_DAYS** — TTL for undelivered (default 14)

* **ALLOWED_ORIGINS** — comma-separated allowed Origins (empty = any)

* **JANITOR_INTERVAL_SECONDS** — cleanup interval (default 60)

* **READ_TIMEOUT_SECONDS / WRITE_TIMEOUT_SECONDS / HANDSHAKE_TIMEOUT_SECONDS** — WS/HTTP timeouts


**Sample `.env.example:`**

``` bash
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
```
---


### API(MVP)

* `GET /healthz → 200 ok`
* `WS /ws?user_id=<id>` — main channel

**Frames (JSON, binary fields are base64):**

```json
// client -> server
{ "type":"send", "to":"<user>", "message_id":"<uuid>", "header_b64":"...", "nonce_b64":"...", "cipher_b64":"..." }
{ "type":"ack",  "message_id":"<uuid>" }
{ "type":"ping" }

// server -> client
{ "type":"deliver", "to":"<user>", "message_id":"<uuid>", "header_b64":"...", "nonce_b64":"...", "cipher_b64":"..." }
{ "type":"pong" }
```
---


### Testing

* Unit tests (Go): `go test ./...` (inside `server/`)
* Integration:
   * 1. `docker compose up -d`
   * 2. Connect WS client to `ws://localhost:8080/ws?user_id=<id>`
   * 3. Send `send`, recive `deliver`, respond with ack

---

### Roadmap

* Protocol: **X3DH + Double Ratchet** (Signal-style)

* **Sealed sender** (hide sender from server)

* Horizontal fan-out: **Redis/NATS** pub/sub

* DB partitioning/sharding; **UUIDv7/ULID** ids

* Attachments via **S3/MinIO** with client-side chunk encryption (AEAD)

* Observability: **Prometheus**, /readyz, /livez

* Alternative transports (gRPC streams), optional Tor/.onion

* Desktop (Tauri) and Mobile (Flutter) clients


---

### Contacts
* Authors/team: 
   riabinkinroman826@gmail.com 
   xvhgxvhg6@gmail.com
* Issues/feature requests: GitHub Issues



---

### License

TBD — MIT 



# AntiPidorandel: Приватный E2EE-мессенджер

## 🇷🇺 Русский

> **Ветка по умолчанию**: `master` (весь актуальный код и документация находятся в этой ветке).

**AntiPidorandel** — приватный мессенджер. Сервер выступает как **«слепой» ретранслятор**: он хранит и пересылает только **зашифрованные конверты** (E2EE). Контент никогда не расшифровывается на сервере, а метаданные сведены к минимуму.

### Документация

Ссылка: https://romanryabinkin.github.io/AntiPidorandel/ru/


---

### Возможности

* **Сервер-ретранслятор (Go)**
  * WebSocket API: `WS /ws?user_id=<id>`
  * Мгновенная доставка онлайн-адресатам, отложенная (очередь) — офлайн-адресатам
  * Подтверждения доставки (`ack`)
* **Зашифрованное хранение (PostgreSQL)**
  * Хранит только шифртекст + минимальные данные маршрутизации (`to_user`, метки времени, TTL)
  * Автоматическая очистка доставленных и просроченных сообщений
* **Эксплуатация и запуск**
  * Health-проба: `GET /healthz`
  * **Docker** / **docker-compose** (distroless-рантайм)
  * Конфигурация через переменные окружения
* **Модель безопасности**
  * E2EE по дизайну (сервер не видит открытый текст)
  * Минимизация логов и метаданных

---

### Структура проекта

* **server/** — самостоятельный модуль Go (сервер)
  * `cmd/server/` — точка входа
  * `internal/transport/ws/` — WebSocket-транспорт
  * `internal/store/` — интерфейс хранилища
    * `postgres/` — реализация на pgx
  * `internal/hub/` — учёт онлайн-сессий
  * `internal/config/` — конфигурация из ENV
  * `internal/wire/` — контракт JSON-кадров «на проводе»
  * `DockerFile` — многостадийная сборка (distroless-рантайм)
  * `go.mod`, `go.sum`, `.gitignore`, `.dockerignore`
* **docker-compose.yml** (в корне репозитория) — локальный запуск Postgres + сервера
* **go.work** (в корне) — Go-workspace, подключает `./server`

---

### Технологии и библиотеки

* Go 1.23+ (`gorilla/websocket`, `pgx`)
* PostgreSQL 16+
* Docker / docker-compose
* Distroless-рантайм (не root)
* Рекомендуемая IDE: VS Code (Go, Docker, YAML, SQLTools, Error Lens)

---

### Быстрый старт

**Требования:**
* Docker Desktop / Docker Engine
* Свободные порты `8080` (сервер) и `5432` (Postgres)

**Запуск через docker-compose (рекомендуется):**
```bash
docker compose up -d --build
docker compose logs -f server
curl http://localhost:8080/healthz   # -> ok
```

**Ручной запуск (без Docker'а):**
```bash
# Start Postgres and create DB "anti"
export DATABASE_URL="postgres://anti:anti@localhost:5432/anti?sslmode=disable"
cd server
go run ./cmd/server
```

---

### Конфигурация

**Переменные окружения (Сервер):**

* **DATABASE_URL** — строка подключения к Postgres (обязательно)

* **ADDR** — адрес прослушивания HTTP/WS (по умолчанию :8080)

* **PENDING_BATCH** — сколько «недоставленных» сообщений отдавать за один проход (по умолчанию 1000)

* **RETAIN_DELIVERED_MINUTES** — TTL для доставленных сообщений (по умолчанию 60 минут)

* **RETAIN_UNDELIVERED_DAYS** — TTL для недоставленных сообщений (по умолчанию 14 дней)

* **ALLOWED_ORIGINS** — разрешённые Origin через запятую (пусто = любые)

* **JANITOR_INTERVAL_SECONDS** —  интервал фоновой очистки (по умолчанию 60 сек)

* **READ_TIMEOUT_SECONDS / WRITE_TIMEOUT_SECONDS / HANDSHAKE_TIMEOUT_SECONDS** — таймауты WS/HTTP


**Пример `.env.example:`**

``` bash
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
```
---


### API(MVP)

* `GET /healthz → 200 ok`
* `WS /ws?user_id=<id>` — основной канал

**Кадры (JSON, бинарные поля — base64):**

```json
// клиент -> сервер
{ "type":"send", "to":"<user>", "message_id":"<uuid>", "header_b64":"...", "nonce_b64":"...", "cipher_b64":"..." }
{ "type":"ack",  "message_id":"<uuid>" }
{ "type":"ping" }

// сервер -> клиент
{ "type":"deliver", "to":"<user>", "message_id":"<uuid>", "header_b64":"...", "nonce_b64":"...", "cipher_b64":"..." }
{ "type":"pong" }
```
---


### Тестирование

* Юнит-тесты (Go): `go test ./...` (внутри `server/`)
* Интеграционные:
   * 1. `docker compose up -d`
   * 2. Подключите WS клиент к `ws://localhost:8080/ws?user_id=<id>`
   * 3. Отправьте `send`, получите `deliver`, ответьте `ack`

---

### Дорожная карта

* Протокол: **X3DH + Double Ratchet** (В стиле Signal)

* **Sealed sender** (скрыть отправителя от серверва)

* Горизонтальный fan-out: **Redis/NATS** pub/sub

* Партиционирование/шардирование БД; идентификаторы **UUIDv7/ULID**

* Вложения **S3/MinIO** с клиентским шифрованием чанков (AEAD)

* Наблюдаемость: **Prometheus**, /readyz, /livez

* Альтернативные транспорты (gRPC streams), опционально Tor/.onion

* Клиенты: Desktop (Tauri) и Mobile (Flutter)


---

### Контакты
* Авторы/Команда: 
   riabinkinroman826@gmail.com 
   xvhgxvhg6@gmail.com
* Issues/запросы фич: GitHub Issues



---

### License

TBD — MIT 








