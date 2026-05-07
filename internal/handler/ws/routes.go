package handler

import (
	"gochat/internal/protocol"
	"gochat/internal/ws"
)

func init() {
	ws.Default.Register(protocol.Login, &ws.Route{SessionFree: true, Handler: HandleLogin})
	ws.Default.Register(protocol.Resume, &ws.Route{SessionFree: true, Handler: Resume})
	ws.Default.Register(protocol.Ping, &ws.Route{Handler: Ping})
	ws.Default.Register(protocol.EnterGroup, &ws.Route{Handler: EnterGroup})
	ws.Default.Register(protocol.EnterLobby, &ws.Route{Handler: EnterLobby})
	ws.Default.Register(protocol.EnterMyGroup, &ws.Route{Handler: EnterMyGroup})
	ws.Default.Register(protocol.SendChat, &ws.Route{Handler: SendChat})
}
