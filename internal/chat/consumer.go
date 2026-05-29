package chat

import (
	"encoding/json"
	"strconv"
	"strings"

	infranats "gochat/internal/infra/nats"
	"gochat/internal/protocol"
	"gochat/internal/ws"
)

func StartGroupChatConsumer(serverName string) error {
	nc, err := infranats.GetNats()
	if err != nil {
		return err
	}
	return nc.SubscribeGroupChat(serverName, func(subject string, data []byte) {
		parts := strings.Split(subject, ".")
		if len(parts) != 3 {
			return
		}
		gid, err := strconv.Atoi(parts[2])
		if err != nil {
			return
		}
		wsMsg, err := json.Marshal(&protocol.Payload{MsgType: protocol.CastChat, Data: data})
		if err != nil {
			return
		}
		ws.BroadcastToGroup(gid, wsMsg)
	})
}
