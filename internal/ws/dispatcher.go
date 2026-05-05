package ws

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/gorilla/websocket"
	"gochat/internal/protocol"
	"sync"
)

type Ctx struct {
	ConnID  string
	Conn    *websocket.Conn
	Payload protocol.Payload
}

type HandlerFunc func(ctx *Ctx) ([]byte, error)

type Route struct {
	SessionFree bool
	Handler     HandlerFunc
}

type Dispatcher struct {
	routes map[protocol.Type]*Route
	mu     sync.RWMutex
}

var Default = NewDispatcher()

func NewDispatcher() *Dispatcher {
	return &Dispatcher{
		routes: map[protocol.Type]*Route{},
	}
}

func (d *Dispatcher) Register(msgType protocol.Type, route *Route) {
	d.mu.Lock()
	d.routes[msgType] = route
	d.mu.Unlock()
}

func (d *Dispatcher) Dispatch(client *Client, in []byte) ([]byte, error) {
	in = bytesTrimSpace(in)

	var p protocol.Payload
	if len(in) == 0 || json.Unmarshal(in, &p) != nil || p.MsgType == "" {
		_ = writeJSON(client.Conn, protocol.Payload{MsgType: protocol.Error, Remark: "invalid payload"})
		_ = client.Conn.Close()
		return nil, errors.New("invalid payload")
	}

	d.mu.RLock()
	h, ok := d.routes[p.MsgType]
	d.mu.RUnlock()

	if !ok {
		_ = writeJSON(client.Conn, protocol.Payload{MsgType: protocol.Error, Remark: "unknown msg type"})
		_ = client.Conn.Close()
		return nil, errors.New("unknown msg type")
	}

	ctx := &Ctx{ConnID: client.ConnID, Conn: client.Conn, Payload: p}

	if !h.SessionFree {
		if GetSocketSession(ctx.ConnID) == nil {
			_ = writeJSON(client.Conn, protocol.Payload{MsgType: protocol.Error, Remark: "session lost"})
			_ = client.Conn.Close()
			return nil, errors.New("session lost")
		}
	}

	return h.Handler(ctx)
}

func writeJSON(conn *websocket.Conn, v any) error {
	b, err := json.Marshal(v)
	if err != nil {
		return err
	}
	return conn.WriteMessage(websocket.TextMessage, b)
}

func bytesTrimSpace(b []byte) []byte {
	return bytes.TrimSpace(b)
}
