package ws

import (
	"crypto/rand"
	"fmt"
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
	PongWait     = 60 * time.Second
	writeTimeout = 10 * time.Second
)

type Client struct {
	ConnID   string
	Conn     *websocket.Conn
	Send     chan []byte
	UserId   int
	ApiToken string
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
	//ticker 搭配內建SetPongHandler維持連線生命
	ticker := time.NewTicker(pingInterval)
	defer func() {
		ticker.Stop()
		Unregister(c.ConnID) //
		c.Conn.Close()
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
			r := redis.GetRedis()
			_ = r.Expire(protocol.SessionKey(c.UserId, c.ConnID), protocol.SessionTTL)
			if c.ApiToken != "" {
				_ = r.Expire(protocol.ApiTokenKey(c.ApiToken), protocol.ApiTokenTTL)
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
	defer hub.mu.Unlock()
	hub.clients[c.ConnID] = c
}

func Unregister(connID string) {
	hub.mu.Lock()
	c := hub.clients[connID]
	delete(hub.clients, connID)
	hub.mu.Unlock()

	ResetGroupScene(connID)
	if c != nil {
		r := redis.GetRedis()
		_ = r.Del(protocol.SessionKey(c.UserId, connID))
		if c.ApiToken != "" {
			_ = r.Del(protocol.ApiTokenKey(c.ApiToken))
		}
	}
}

// groupSubs tracks which clients are currently viewing each group on one server.
var groupSubs = struct {
	byGroup map[int]map[string]*Client
	byConn  map[string]int
	mu      sync.RWMutex
}{
	byGroup: make(map[int]map[string]*Client),
	byConn:  make(map[string]int),
}

func JoinGroup(connID string, groupID int, c *Client) {
	groupSubs.mu.Lock()
	defer groupSubs.mu.Unlock()

	if oldID, ok := groupSubs.byConn[connID]; ok {
		delete(groupSubs.byGroup[oldID], connID)
	}
	if groupSubs.byGroup[groupID] == nil {
		groupSubs.byGroup[groupID] = make(map[string]*Client)
	}
	groupSubs.byGroup[groupID][connID] = c
	groupSubs.byConn[connID] = groupID
}

func ResetGroupScene(connID string) {
	groupSubs.mu.Lock()
	defer groupSubs.mu.Unlock()

	if groupID, ok := groupSubs.byConn[connID]; ok {
		delete(groupSubs.byGroup[groupID], connID)
		delete(groupSubs.byConn, connID)
	}
}

func BroadcastToGroup(groupID int, data []byte) {
	groupSubs.mu.RLock()
	defer groupSubs.mu.RUnlock()

	for _, c := range groupSubs.byGroup[groupID] {
		c.TrySend(data)
	}
}
