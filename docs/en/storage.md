# Storage (PostgreSQL)

## Table `messages`
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

- Delivered messages are deleted after `RETAIN_DELIVERED_MINUTES`.
- Undelivered messages — after `RETAIN_UNDELIVERED_DAYS`.
- Mailbox nodes can replicate via `NODE_PEERS` (optional).
