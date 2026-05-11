package ws

import (
	"crypto/rand"
	"fmt"
	log2 "gochat/internal/log"
	"log"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	"gochat/internal/infra/redis"
	"gochat/internal/protocol"
)

const (
	sendBufSize  = 64
	pingInterval = 30 * time.Second
	writeTimeout = 10 * time.Second
)

type Client struct {
	ConnID string
	Conn   *websocket.Conn
	Send   chan []byte
}

func NewClient(conn *websocket.Conn) *Client {
	b := make([]byte, 8)
	_, _ = rand.Read(b)
	return &Client{
		ConnID: fmt.Sprintf("%x", b),
		Conn:   conn,
		Send:   make(chan []byte, sendBufSize),
	}
}

func (c *Client) WritePump() {
	//ticker 搭配內建SetPongHandler
	ticker := time.NewTicker(pingInterval)
	defer func() {
		ticker.Stop()
		Unregister(c.ConnID)
		c.Conn.Close()
		_ = log2.SaveLog(log2.TypeConnErr, c.ConnID)
	}()

	for {
		select {
		case msg, ok := <-c.Send:
			_ = c.Conn.SetWriteDeadline(time.Now().Add(writeTimeout))
			if !ok {
				c.Conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}
			if err := c.Conn.WriteMessage(websocket.TextMessage, msg); err != nil {
				log.Printf("[WS] write error: %v", err)
				return
			}
		case <-ticker.C:
			_ = c.Conn.SetWriteDeadline(time.Now().Add(writeTimeout))
			if err := c.Conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				log.Printf("[WS] ping error: %v", err)
				return
			}
		}
	}
}

func (c *Client) TrySend(msg []byte) {
	select {
	case c.Send <- msg:
	default:
	}
}

var hub = struct {
	clients map[string]*Client
	mu      sync.RWMutex
}{clients: make(map[string]*Client)}

func Register(c *Client) {
	hub.mu.Lock()
	hub.clients[c.ConnID] = c
	hub.mu.Unlock()
}

func Unregister(connID string) {
	hub.mu.Lock()
	delete(hub.clients, connID)
	hub.mu.Unlock()

	_ = redis.GetRedis().Del(protocol.SessionKey(connID))
}

func GetClient(connID string) (*Client, bool) {
	hub.mu.RLock()
	c, ok := hub.clients[connID]
	hub.mu.RUnlock()
	return c, ok
}

func Broadcast(msg []byte) {
	hub.mu.RLock()
	defer hub.mu.RUnlock()
	for _, c := range hub.clients {
		c.TrySend(msg)
	}
}
