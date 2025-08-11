package wire

// FrameType задаёт тип кадра протокола, идущего по WebSocket.
type FrameType string


const (
	// TypeSend — клиент отправляет серверу зашифрованный конверт для получателя.
	TypeSend FrameType = "send"      
	// TypeDeliver — сервер доставляет клиенту зашифрованный конверт.
	TypeDeliver FrameType = "deliver"
	// TypeAck — клиент подтверждает доставку сообщения по его MessageID.
	TypeAck FrameType = "ack"      
	// TypePing — Дефолтный клиентский пинг
	TypePing FrameType = "ping"
	// TypePong — ответ сервера на пинг.
	TypePong FrameType = "pong"
)

// Frame представляет универсальный кадр протокола (JSON) для WebSocket соединения.
// Поля *_B64 содержат бинарные данные в кодировке base64.
type Frame struct {
	Type      FrameType `json:"type"`
	To        string    `json:"to,omitempty"`
	MessageID string    `json:"message_id,omitempty"`
	HeaderB64 string    `json:"header_b64,omitempty"`
	NonceB64  string    `json:"nonce_b64,omitempty"`
	CipherB64 string    `json:"cipher_b64,omitempty"`
	Ping      string    `json:"ping,omitempty"`
	Pong      string    `json:"pong,omitempty"`
}