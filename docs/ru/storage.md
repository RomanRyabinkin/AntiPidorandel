# Хранилище (PostgreSQL)

## Таблица `messages`
```sql
CREATE TABLE IF NOT EXISTS messages (
  id           UUID PRIMARY KEY,
  to_user      TEXT        NOT NULL,
  header       BYTEA,
  nonce        BYTEA       NOT NULL,
  cipher       BYTEA       NOT NULL,
  created_at   TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  delivered_at TIMESTAMPTZ,
  expires_at   TIMESTAMPTZ
);
CREATE INDEX IF NOT EXISTS idx_inbox_pending
  ON messages (to_user, created_at) WHERE delivered_at IS NULL;
CREATE INDEX IF NOT EXISTS idx_expiry ON messages (expires_at);
```

- Доставленные удаляются спустя `RETAIN_DELIVERED_MINUTES`.
- Недоставленные — по `RETAIN_UNDELIVERED_DAYS`.
- Узлы‑почтовики могут реплицировать по `NODE_PEERS` (опционально).
