package ws

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/RomanRyabinkin/antipidorandel/internal/config"
	"github.com/RomanRyabinkin/antipidorandel/internal/hub"
	"github.com/RomanRyabinkin/antipidorandel/internal/store"
	"github.com/RomanRyabinkin/antipidorandel/internal/wire"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

// Обслуживает WS-подключения и взаимодействует со Store.
// Не имеет доступа к криптографии и содержимого сообщения.
type Server struct {
	H   *hub.Hub
	St  store.Store
	Cfg config.Config
}

// Создает сервер WebSocket транспорта
func New(h *hub.Hub, st store.Store, cfg config.Config) *Server {
	return &Server{H: h, St: st, Cfg: cfg}
}

// Регистрирует обработчик WebSocket по пути /ws на указанном ServeMux.
//
// Контракт подключения:
//   - Обязательный query-параметр user_id определяет «почтовый ящик» получателя.
//   - Проверка Origin настраивается через cfg.AllowedOrigins (пустой список — любой Origin).
//   - После апгрейда соединения Server:
//       1) регистрирует клиента в hub,
//       2) асинхронно отправляет все pending-сообщения получателю,
//       3) запускает цикл записи (пинг/понг keepalive),
//       4) читает входящие кадры и обрабатывает: send/ack/ping.
//
// Закрытие соединения приводит к deregister клиента и остановке его очереди отправки.
func (s *Server) Router(mux *http.ServeMux) {
	upgrader := websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			if len(s.Cfg.AllowedOrigins) == 0 { return true }
			origin := r.Header.Get("Origin")
			for _, o := range s.Cfg.AllowedOrigins { if o == origin { return true } }
			return false
		},
	}
	mux.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		user := r.URL.Query().Get("user_id")
		if user == "" {
			http.Error(w, "user_id required", http.StatusUnauthorized); return
		}
		ws, err := upgrader.Upgrade(w, r, nil)
		if err != nil { log.Println("ws upgrade:", err); return }

		_ = ws.SetReadDeadline(time.Now().Add(s.Cfg.ReadTimeout))
		ws.SetPongHandler(func(string) error {
			_ = ws.SetReadDeadline(time.Now().Add(s.Cfg.ReadTimeout)); return nil
		})

		c := &hub.Client{User: user, WS: ws, Send: make(chan []byte, 256)}
		s.H.Set(user, c)
		log.Printf("online: %s", user)

		// writer: отправка, keepalive ping
		go func() {
			defer ws.Close()
			ticker := time.NewTicker(20 * time.Second)
			defer ticker.Stop()
			for {
				select {
				case b, ok := <-c.Send:
					_ = ws.SetWriteDeadline(time.Now().Add(s.Cfg.WriteTimeout))
					if !ok {
						_ = ws.WriteMessage(websocket.CloseMessage, []byte{}); return
					}
					if err := ws.WriteMessage(websocket.TextMessage, b); err != nil { return }
				case <-ticker.C:
					_ = ws.WriteControl(websocket.PingMessage, []byte("p"), time.Now().Add(5*time.Second))
				}
			}
		}()

		// отдача недоставленных сообщений
		go s.deliverPending(r.Context(), c)

		// reader: принятие кадров
		for {
			_, data, err := ws.ReadMessage()
			if err != nil { break }
			var f wire.Frame
			if err := json.Unmarshal(data, &f); err != nil { continue }
			switch f.Type {
			case wire.TypeSend:
				s.onSend(r.Context(), f)
			case wire.TypeAck:
				if id, err := uuid.Parse(f.MessageID); err == nil {
					_ = s.St.MarkDelivered(r.Context(), id)
				}
			case wire.TypePing:
				_ = ws.WriteJSON(wire.Frame{Type: wire.TypePong})
			default:
				// игнор
			}
		}
		s.H.Del(user)
		close(c.Send)
		log.Printf("offline: %s", user)
	})
}


// Выгружает из Store все недоставленные сообщения для клиента c
// и по одному пушит их в его исходящую очередь. Если очередь блокируется дольше
// 5 секунд, функция прекращает доставку (неудачные сообщения останутся в Store).
func (s *Server) deliverPending(ctx context.Context, c *hub.Client) {
	pending, err := s.St.ListPending(ctx, c.User, 1000)
	if err != nil { log.Println("list pending:", err); return }
	for _, m := range pending {
		df := wire.Frame{
			Type:      wire.TypeDeliver,
			To:        c.User,
			MessageID: m.ID.String(),
			HeaderB64: b64(m.Header),
			NonceB64:  b64(m.Nonce),
			CipherB64: b64(m.Cipher),
		}
		b, _ := json.Marshal(df)
		select {
		case c.Send <- b:
		case <-time.After(5 * time.Second):
			log.Println("deliver queue timeout"); return
		}
	}
}

// Обрабатывает входящий кадр send:
//   1) валидирует идентификатор сообщения;
//   2) декодирует бинарные поля из base64;
//   3) сохраняет конверт в Store с TTL для недоставленных;
//   4) если получатель онлайн — немедленно отправляет deliver.
//
// Функция умышленно «молчит» об ошибках (логирует и продолжает),
// чтобы не раскрывать серверные детали/состояния.
func (s *Server) onSend(ctx context.Context, f wire.Frame) {
	msgID, err := uuid.Parse(f.MessageID); if err != nil { return }
	header, _ := b64d(f.HeaderB64)
	nonce,  _ := b64d(f.NonceB64)
	ciph,   _ := b64d(f.CipherB64)
	exp := time.Now().Add(s.Cfg.RetainUndelivered)

	_ = s.St.Save(ctx, store.Message{
		ID: msgID, ToUser: f.To, Header: header, Nonce: nonce, Cipher: ciph, ExpiresAt: &exp,
	})

	// если получатель онлайн — доставка осуществляется сразу
	if rc, ok := s.H.Get(f.To); ok {
		df := wire.Frame{Type: wire.TypeDeliver, To: f.To, MessageID: f.MessageID, HeaderB64: f.HeaderB64, NonceB64: f.NonceB64, CipherB64: f.CipherB64}
		if b, err := json.Marshal(df); err == nil {
			select {
			case rc.Send <- b:
			default:
				log.Printf("deliver backlog for %s", f.To)
			}
		}
	}
}

// Кодирует байты в base64 для передачи в JSON
func b64(b []byte) string {
	if len(b) == 0 { return "" }
	return base64.StdEncoding.EncodeToString(b)
}

// Декодирует base64-строку в байты
// При пустой строке возвращает (nil, nil)
func b64d(s string) ([]byte, error) {
	if s == "" { return nil, nil }
	return base64.StdEncoding.DecodeString(s)
}
