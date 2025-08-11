package postgres

import (
	"context"
	"time"

	"github.com/RomanRyabinkin/antipidorandel/internal/store"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Store struct{ db *pgxpool.Pool }

// Создает соединение с постгрой(PostgreSQL) по DATABASE_URL
func New(ctx context.Context, url string) (*Store, error) {
	cfg, err := pgxpool.ParseConfig(url);
	if err != nil {
		return nil, err
	}
	db, err := pgxpool.NewWithConfig(ctx, cfg)
	if err != nil {
		return nil, err
	}
	return &Store{db: db}, nil
}

// Закрывает пул коннектов(соединений)
func (s *Store) Close() { s.db.Close() }

// Создает таблицы и индексы, если их еще нет
func (s *Store) Init(ctx context.Context) error {
	_, err := s.db.Exec(
		ctx,
		`
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
`,
	)
	return err
}

// Сохраняет зашифрованный конверт
// Повторная вставка по тому же ID игнорится
func (s * Store) Save(ctx context.Context, m store.Message) error {
	_, err := s.db.Exec(ctx, `
	INSERT INTO messages (id, to_user, header, nonce, cipher, created_at, expires_at)
VALUES ($1,$2,$3,$4,$5,NOW(),$6)
ON CONFLICT (id) DO NOTHING
	`, m.ID, m.ToUser, m.Header, m.Nonce, m.Cipher, m.ExpiresAt)
	return err
}

// Возвращает недоставленные сообщения для пользователя
func (s * Store) ListPending(ctx context.Context, toUser string, limit int) ([]store.Message, error) {
	rows, err := s.db.Query(ctx, `
	SELECT id, to_user, header, nonce, cipher, created_at, delivered_at, expires_at
FROM messages
WHERE to_user=$1 AND delivered_at IS NULL
ORDER BY created_at ASC
LIMIT $2
	`, toUser, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []store.Message
	for rows.Next() {
		var m store.Message
		if err := rows.Scan(&m.ID, &m.ToUser, &m.Header, &m.Nonce, &m.Cipher, &m.CreatedAt, &m.DeliveredAt, &m.ExpiresAt); err != nil {
			return nil, err
		}
		out = append(out, m)
	}
	return out, rows.Err()
}

// Помечает сообщение доставленным по ID
func (s *Store) MarkDelivered(ctx context.Context, id uuid.UUID) error {
	_, err := s.db.Exec(ctx, `UPDATE messages SET delivered_at=NOW() WHERE id=$1`, id)
	return err
}

type UUID = [16]byte

// Удаляет доставленные сообщения раньше deliveredBefore и просроченные на момент текущей UTC метки
func (s *Store) DeleteExpired(ctx context.Context, deliveredBefore time.Time, now time.Time) (int64, int64, error) {
	ct1, err := s.db.Exec(ctx, `
DELETE FROM messages
 WHERE delivered_at IS NOT NULL
   AND delivered_at < $1`, deliveredBefore)
	if err != nil {
		return 0, 0, err
	}
	ct2, err := s.db.Exec(ctx, `
DELETE FROM messages
 WHERE expires_at IS NOT NULL
   AND expires_at < $1`, now)
	if err != nil {
		return ct1.RowsAffected(), 0, err
	}
	return ct1.RowsAffected(), ct2.RowsAffected(), nil
}