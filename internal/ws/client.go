package ws

import (
	"crypto/rand"
	"fmt"
	"github.com/gorilla/websocket"
	"sync"
)

type Client struct {
	ConnID string
	Conn   *websocket.Conn
}

func NewClient(conn *websocket.Conn) *Client {
	b := make([]byte, 8)
	_, _ = rand.Read(b)
	return &Client{
		ConnID: fmt.Sprintf("%x", b),
		Conn:   conn,
	}
}

type Hub struct {
	clients map[string]*Client
	mu      sync.RWMutex
}

var defaultHub = &Hub{
	clients: make(map[string]*Client),
}

func (h *Hub) register(c *Client) {
	h.mu.Lock()
	h.clients[c.ConnID] = c
	h.mu.Unlock()
}

func (h *Hub) unregister(connID string) {
	h.mu.Lock()
	delete(h.clients, connID)
	h.mu.Unlock()
}

func (h *Hub) get(connID string) (*Client, bool) {
	h.mu.RLock()
	c, ok := h.clients[connID]
	h.mu.RUnlock()
	return c, ok
}

func Register(c *Client)                       { defaultHub.register(c) }
func Unregister(connID string)                 { defaultHub.unregister(connID) }
func GetClient(connID string) (*Client, bool)  { return defaultHub.get(connID) }
