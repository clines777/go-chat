package consumer

import (
	"encoding/json"
	gonats "github.com/nats-io/nats.go"
	"gochat/internal/protocol"
	"gochat/internal/ws"
	"log"
	"strconv"
	"strings"
)

func CastLeaveGroup(msg gonats.Msg) {

}

// CastUpdateGroup - 轉發群組資料更新, 純通知, AckNone
func CastUpdateGroup(msg *gonats.Msg) {
	parts := strings.Split(msg.Subject, ".")
	if len(parts) != 3 {
		return
	}
	gid, err := strconv.Atoi(parts[2])
	if err != nil {
		return
	}
	castMsg, err := json.Marshal(&protocol.Payload{MsgType: protocol.UpdateGroupCast, Data: msg.Data})
	if err != nil {
		log.Printf("group update consumer marshal error: %v", err)
		return
	}
	ws.BroadcastToGroup(gid, castMsg)
}

// CastSendChat - 轉發最新訊息給群內, 純通知, AckNone
func CastSendChat(msg *gonats.Msg) {
	parts := strings.Split(msg.Subject, ".")
	if len(parts) != 3 {
		return
	}
	gid, err := strconv.Atoi(parts[2])
	if err != nil {
		return
	}
	wsMsg, err := json.Marshal(&protocol.Payload{MsgType: protocol.CastChat, Data: msg.Data})
	if err != nil {
		log.Printf("group chat consumer error: %v", err)
		return
	}

	ws.BroadcastToGroup(gid, wsMsg)
}
