package consumer

import (
	"encoding/json"
	"log"
	"strconv"
	"strings"

	"gochat/internal/protocol"
	"gochat/internal/ws"
)

func GroupChat(subject string, data []byte) {
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
		log.Printf("group chat consumer error: %v", err)
		return
	}

	ws.BroadcastToGroup(gid, wsMsg)
}
