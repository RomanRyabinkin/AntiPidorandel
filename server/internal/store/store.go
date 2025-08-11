package store

import (
	"context"
	"time"

	"github.com/google/uuid"
)


// Message описывает зашифрованный «конверт», который сервер хранит и доставляет.
type Message struct {
	ID          uuid.UUID   // генерируется клиентом (лучше UUIDv7/ULID)
	ToUser      string      // только получатель; отправителя сервер не хранит
	Header      []byte      // опциональный зашифрованный заголовок
	Nonce       []byte      // nonce для AEAD
	Cipher      []byte      // шифр-текст(включая MAC)
	CreatedAt   time.Time
	DeliveredAt *time.Time
	ExpiresAt   *time.Time  // TTL для недоставленных сообщений
}

// Store описывает минимальные операции хранилища сообщений.
type Store interface {
	// Создает необходимые таблицы/индексы
	Init(ctx context.Context) error
	// Сохраняет сообщение (будет идемпотентным по ID)
	Save(ctx context.Context, m Message) error
	// Возвращает недоставленные сообщения для получателя
	ListPending(ctx context.Context, toUser string, limit int) ([]Message, error)
	// Помечает сообщение доставленным по ID
	MarkDelivered(ctx context.Context, id uuid.UUID) error
	// Удаляет устаревшие записи: доставленные раньше deliveredBefore и просроченные на текущую UTC метку
	DeleteExpired(ctx context.Context, deliveredBefore time.Time, now time.Time) (deliveredRemoved int64, expiredRemoved int64, err error)
	// Освобождает ресурсы хранилища
	Close()
}