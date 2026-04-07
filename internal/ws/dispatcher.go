package ws

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/gorilla/websocket"
	handler "gochat/internal/handler/ws"
	"gochat/internal/protocol"
	"sync"
)

type Ctx struct {
	Conn    *websocket.Conn
	Payload protocol.Payload
}

type HandlerFunc func(ctx *Ctx) ([]byte, error)

type Route struct {
	SessionFree bool
	Handler     HandlerFunc
}

type Dispatcher struct {
	routes           map[protocol.Type]*Route
	checkSessionFree map[protocol.Type]bool
	mu               sync.RWMutex
}

func NewDispatcher() *Dispatcher {
	dispatcher := &Dispatcher{
		routes:           map[protocol.Type]*Route{},
		checkSessionFree: map[protocol.Type]bool{},
	}

	dispatcher.checkSessionFree = make(map[protocol.Type]bool)
	dispatcher.checkSessionFree[protocol.Login] = true
	dispatcher.checkSessionFree[protocol.Resume] = true

	dispatcher.routes[protocol.Login] = &Route{
		SessionFree: true,
		Handler:     handler.Login,
	}

	dispatcher.routes[protocol.Resume] = &Route{
		SessionFree: true,
		Handler:     handler.Resume,
	}

	dispatcher.routes[protocol.Ping] = &Route{
		Handler: handler.Ping,
	}

	dispatcher.routes[protocol.EnterGroup] = &Route{
		Handler: handler.EnterGroup,
	}

	dispatcher.routes[protocol.EnterLobby] = &Route{
		Handler: handler.EnterLobby,
	}

	dispatcher.routes[protocol.EnterMyGroup] = &Route{
		Handler: handler.EnterMyGroup,
	}

	dispatcher.routes[protocol.SendChat] = &Route{
		Handler: handler.SendChat,
	}

	return dispatcher
}

func (d *Dispatcher) Dispatch(conn *websocket.Conn, in []byte) ([]byte, error) {
	in = bytesTrimSpace(in)
	var p protocol.Payload
	if len(in) == 0 || json.Unmarshal(in, &p) != nil || p.MsgType == "" {
		_ = writeJSON(conn, protocol.Payload{MsgType: "error", Remark: "invalid payload"})
		_ = conn.Close()
		return nil, errors.New("invalid payload")
	}

	d.mu.RLock()
	h, ok := d.routes[p.MsgType]
	d.mu.RUnlock()
	if !ok {
		_ = writeJSON(conn, protocol.Payload{MsgType: "error", Remark: "unknown id"})
		_ = conn.Close()
		return nil, errors.New("unknown msg type")
	}

	ctx := &Ctx{Conn: conn}

	if !d.checkSessionFree[p.MsgType] {
		sess := GetSocketSession(ctx.Conn)
		if sess == nil {
			_ = writeJSON(conn, protocol.Payload{MsgType: "error", Remark: "session lost"})
			_ = conn.Close()
			return nil, errors.New("session lost")
		}
	}

	out, err := h.Handle(ctx, p)
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
