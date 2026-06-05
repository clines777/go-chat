package ws

import (
	"crypto/rand"
	"errors"
	"fmt"
	"log"
	"net"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	"gochat/internal/infra"
	"gochat/internal/infra/redis"
	"gochat/internal/protocol"
)

const (
	sendBufSize  = 64
	pingInterval = 30 * time.Second
	PongWait     = 60 * time.Second
	writeTimeout = 10 * time.Second
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     checkOrigin,
}

// checkOrigin WS CheckOrigin, 防止跨站。
func checkOrigin(r *http.Request) bool {
	origin := r.Header.Get("Origin")
	if origin == "" {
		return true
	}
	u, err := url.Parse(origin)
	if err != nil || u.Hostname() == "" {
		return false
	}

	if allowed := infra.GetEnvConfig().WSAllowedOrigins; len(allowed) > 0 {
		for _, a := range allowed {
			if strings.EqualFold(origin, strings.TrimSpace(a)) {
				return true
			}
		}
		return false
	}

	reqHost := r.Host
	if h, _, splitErr := net.SplitHostPort(reqHost); splitErr == nil {
		reqHost = h
	}
	return strings.EqualFold(u.Hostname(), reqHost)
}

// Client - 連線物件
type Client struct {
	ConnID   string
	Conn     *websocket.Conn
	Send     chan []byte
	UserId   int
	ApiToken string
}

// HandleWs - 處理ws message
func HandleWs(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("[WS] upgrade error: %v", err)
		return
	}

	client := NewClient(conn)
	go client.WritePump()

	// 讀迴圈離開時統一清理: 註銷連線/清 session → 關閉 Send 通知 WritePump 收尾 → 關連線。
	defer func() {
		Unregister(client.ConnID)
		close(client.Send)
		client.Conn.Close()
	}()

	conn.SetPongHandler(func(string) error {
		_ = conn.SetReadDeadline(time.Now().Add(PongWait))
		return nil
	})

	for {
		mt, input, err := conn.ReadMessage()
		if err != nil {
			log.Printf("[WS] read error: %v", err)
			return
		}
		if mt != websocket.TextMessage {
			return
		}

		output, dispatchErr := Default.Dispatch(client, input)
		if dispatchErr != nil {
			var de *protocol.DispatchError
			if errors.As(dispatchErr, &de) && de.Fatal {
				return
			}
			continue
		}

		if len(output) > 0 {
			select {
			case client.Send <- output:
			default:
				// Send chan滿, 視為slow client, 連線由defer統一清理。
				return
			}
		}
	}
}

// NewClient 新建一個ws連線物件
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
	// 只關連線喚醒讀迴圈, 註銷/清 session 統一由 HandleWs 的 defer 負責, 避免重複清理。
	defer func() {
		ticker.Stop()
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

// TrySend - 推送ws消息
func (c *Client) TrySend(msg []byte) {
	select {
	case c.Send <- msg:
	default:
	}
}

// ws連線hub
var hub = struct {
	clients map[string]*Client
	mu      sync.RWMutex
}{clients: make(map[string]*Client)}

// Register - 註冊ws連線
func Register(c *Client) {
	hub.mu.Lock()
	defer hub.mu.Unlock()
	hub.clients[c.ConnID] = c
}

// Unregister - 註銷ws連線並清除session狀態
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

// groupSubs 將ws連線物件byConn跟byGroup儲存, 便於廣播推送
var groupSubs = struct {
	byGroup map[int]map[string]*Client
	byConn  map[string]int
	mu      sync.RWMutex
}{
	byGroup: make(map[int]map[string]*Client),
	byConn:  make(map[string]int),
}

// JoinGroup - 用戶進入群內場景時將其連線置入groupSubs
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

// ResetGroupScene 重置用戶連線所在場景
func ResetGroupScene(connID string) {
	groupSubs.mu.Lock()
	defer groupSubs.mu.Unlock()

	if groupID, ok := groupSubs.byConn[connID]; ok {
		delete(groupSubs.byGroup[groupID], connID)
		delete(groupSubs.byConn, connID)
	}
}

// BroadcastToGroup 群內廣播
func BroadcastToGroup(groupID int, data []byte) {
	groupSubs.mu.RLock()
	defer groupSubs.mu.RUnlock()

	for _, c := range groupSubs.byGroup[groupID] {
		c.TrySend(data)
	}
}

// SendToGroupUser 發送資料給特定 user 目前正處於該群場景的所有連線。
func SendToGroupUser(groupID, userID int, data []byte) {
	groupSubs.mu.RLock()
	defer groupSubs.mu.RUnlock()

	for _, c := range groupSubs.byGroup[groupID] {
		if c.UserId == userID {
			c.TrySend(data)
		}
	}
}

// EvictGroupUser 將某 user 的連線移出該群場景, 回傳被移出的 connID, 之後群內廣播不再送達。
func EvictGroupUser(groupID, userID int) []string {
	groupSubs.mu.Lock()
	defer groupSubs.mu.Unlock()

	var conns []string
	for connID, c := range groupSubs.byGroup[groupID] {
		if c.UserId == userID {
			delete(groupSubs.byGroup[groupID], connID)
			delete(groupSubs.byConn, connID)
			conns = append(conns, connID)
		}
	}
	return conns
}
