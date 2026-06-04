package handler

import (
	"gochat/internal/protocol"
	"gochat/internal/ws"
)

func init() {
	ws.Default.Register(protocol.Login, &ws.Route{SessionFree: true, Handler: Login})
	ws.Default.Register(protocol.EnterGroup, &ws.Route{Handler: EnterGroup})
	ws.Default.Register(protocol.EnterLobby, &ws.Route{Handler: EnterLobby})
	ws.Default.Register(protocol.EnterMyGroup, &ws.Route{Handler: EnterMyGroup})
	ws.Default.Register(protocol.EnterSelf, &ws.Route{Handler: EnterSelf})
	ws.Default.Register(protocol.SendChat, &ws.Route{Handler: SendChat})
	ws.Default.Register(protocol.Resume, &ws.Route{SessionFree: true, Handler: Resume})
	ws.Default.Register(protocol.JoinGroup, &ws.Route{Handler: JoinGroup})
	ws.Default.Register(protocol.LeaveGroup, &ws.Route{Handler: LeaveGroup})
	ws.Default.Register(protocol.Logout, &ws.Route{Handler: Logout})
	ws.Default.Register(protocol.UpdateLastRead, &ws.Route{Handler: UpdateLastRead})
	ws.Default.Register(protocol.UpdateGroup, &ws.Route{Handler: UpdateGroup})
	ws.Default.Register(protocol.PinChat, &ws.Route{Handler: PinChat})
	ws.Default.Register(protocol.UnpinChat, &ws.Route{Handler: UnpinChat})
	ws.Default.Register(protocol.DelChat, &ws.Route{Handler: DelChat})
	ws.Default.Register(protocol.BanUser, &ws.Route{Handler: BanUser})
	ws.Default.Register(protocol.UnbanUser, &ws.Route{Handler: UnbanUser})
	ws.Default.Register(protocol.KickUser, &ws.Route{Handler: KickUser})
}
