package hub

import (
	"sync"

	"github.com/gorilla/websocket"
)

// Описывает активное соединение пользователя.
type Client struct {
	User string
	WS *websocket.Conn
	Send chan []byte // исходящая очередь сообщений
}

// Хранит соответствие user_id -> Client
type Hub struct {
	clients sync.Map
}

// Создает пустой хаб коннектов(подключений)
func New() *Hub { return &Hub{} }

// Регистрирует/обновляет соединение для пользователя
func (h *Hub) Set(user string, c *Client) { h.clients.Store(user, c) }

// Удаляет соединение пользователя
func (h *Hub) Del(user string)            { h.clients.Delete(user) }

// Возвращает активное соединение пользователя, если оно есть
func (h *Hub) Get(user string) (*Client, bool) {
	v, ok := h.clients.Load(user)
	if !ok { return nil, false }
	return v.(*Client), true
}
