package ws

import (
	"bytes"
	"encoding/json"
	"errors"
	"gochat/internal/protocol"
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

	var p *protocol.Payload
	if len(in) == 0 || json.Unmarshal(in, &p) != nil || p.MsgType == "" {
		writeError(client, "invalid payload")
		return nil, errors.New("invalid payload")
	}

	d.mu.RLock()
	h, ok := d.routes[p.MsgType]
	d.mu.RUnlock()

	if !ok {
		writeError(client, "unknown msg type")
		return nil, errors.New("unknown msg type")
	}

	ctx := &Ctx{Client: client, Payload: p}

	if !h.SessionFree {
		if GetSocketSession(client.ConnID) == nil {
			writeError(client, "session lost")
			return nil, errors.New("session lost")
		}
	}

	out := h.Handler(ctx)
	if out == nil {
		return nil, nil
	}
	b, err := json.Marshal(out)
	if err != nil {
		writeError(client, "internal error")
		return nil, err
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
