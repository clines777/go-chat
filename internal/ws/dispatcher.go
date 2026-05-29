package ws

import (
	"bytes"
	"encoding/json"
	"gochat/internal/protocol"
	"gochat/internal/session"
	"sync"
)

type Ctx struct {
	Client  *Client
	Payload *protocol.Payload
}

type HandlerFunc func(ctx *Ctx) *protocol.Payload

type Route struct {
	SessionFree bool
	Handler     HandlerFunc
}

type Dispatcher struct {
	routes map[protocol.MsgType]*Route
	mu     sync.RWMutex
}

var Default = NewDispatcher()

func NewDispatcher() *Dispatcher {
	return &Dispatcher{
		routes: map[protocol.MsgType]*Route{},
	}
}

func (d *Dispatcher) Register(msgType protocol.MsgType, route *Route) {
	d.mu.Lock()
	d.routes[msgType] = route
	d.mu.Unlock()
}

func (d *Dispatcher) Dispatch(client *Client, in []byte) ([]byte, error) {
	in = bytesTrimSpace(in)

	var p *protocol.Payload
	if len(in) == 0 || json.Unmarshal(in, &p) != nil || p.MsgType == "" {
		writeError(client, "invalid payload")
		return nil, &protocol.DispatchError{Code: 4000, Message: "invalid payload"}
	}

	d.mu.RLock()
	h, ok := d.routes[p.MsgType]
	d.mu.RUnlock()

	if !ok {
		writeError(client, "unknown msg type")
		return nil, &protocol.DispatchError{Code: protocol.ErrUnknownMsgType, Message: "unknown msg type"}
	}

	ctx := &Ctx{Client: client, Payload: p}

	if !h.SessionFree {
		if session.Get(client.UserId, client.ConnID) == nil {
			writeError(client, "session required")
			return nil, &protocol.DispatchError{Code: protocol.ErrSessionRequired, Message: "session required"}
		}
	}

	out := h.Handler(ctx)
	if out == nil {
		return nil, nil
	}
	b, err := json.Marshal(out)
	if err != nil {
		writeError(client, "internal error")
		return nil, &protocol.DispatchError{Code: protocol.ErrInternalError, Message: "internal error", Fatal: true}
	}
	return b, nil
}

func writeError(client *Client, remark string) {
	b, err := json.Marshal(protocol.Payload{MsgType: protocol.Error, Remark: remark})
	if err != nil {
		return
	}
	client.TrySend(b)
}

func bytesTrimSpace(b []byte) []byte {
	return bytes.TrimSpace(b)
}
