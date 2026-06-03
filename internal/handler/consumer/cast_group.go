package consumer

import (
	"encoding/json"
	gonats "github.com/nats-io/nats.go"
	"gochat/internal/protocol"
	"gochat/internal/ws"
	"log"
)

func CastLeaveGroup(msg *gonats.Msg) {
	var event protocol.LeaveGroupEvent
	if err := json.Unmarshal(msg.Data, &event); err != nil {
		log.Printf("[CastLeaveGroup] unmarshall err: %v", err)
		return
	}

	castMsg, err := json.Marshal(&protocol.Payload{MsgType: protocol.LeaveGroupCast, Data: msg.Data})
	if err != nil {
		log.Printf("[CastLeaveGroup] consumer marshal error: %v", err)
		return
	}

	ws.BroadcastToGroup(event.GroupId, castMsg)
}

// CastUpdateGroup - 轉發群組資料更新, 純通知, AckNone
func CastUpdateGroup(msg *gonats.Msg) {
	var event protocol.UpdateGroupCastEvent
	if err := json.Unmarshal(msg.Data, &event); err != nil {
		log.Printf("[CastUpdateGroup] unmarshall err: %v", err)
		return
	}
	castMsg, err := json.Marshal(&protocol.Payload{MsgType: protocol.UpdateGroupCast, Data: msg.Data})
	if err != nil {
		log.Printf("[CastUpdateGroup] consumer marshal error: %v", err)
		return
	}
	ws.BroadcastToGroup(event.GroupID, castMsg)
}

// CastSendChat - 轉發最新訊息給群內, 純通知, AckNone
func CastSendChat(msg *gonats.Msg) {
	var event protocol.CastChatEvent
	if err := json.Unmarshal(msg.Data, &event); err != nil {
		log.Printf("[CastSendChat] unmarshall err: %v", err)
		return
	}
	wsMsg, err := json.Marshal(&protocol.Payload{MsgType: protocol.CastChat, Data: msg.Data})
	if err != nil {
		log.Printf("[CastSendChat] consumer error: %v", err)
		return
	}

	ws.BroadcastToGroup(event.GroupID, wsMsg)
}

// CastPinChat - 轉發置頂發言, 純通知, AckNone
func CastPinChat(msg *gonats.Msg) {
	var event protocol.PinChatCastEvent
	if err := json.Unmarshal(msg.Data, &event); err != nil {
		log.Printf("[CastPinChat] unmarshall err: %v", err)
		return
	}
	castMsg, err := json.Marshal(&protocol.Payload{MsgType: protocol.PinChatCast, Data: msg.Data})
	if err != nil {
		log.Printf("[CastPinChat] consumer marshal error: %v", err)
		return
	}
	ws.BroadcastToGroup(event.GroupID, castMsg)
}

// CastUnpinChat - 轉發取消置頂, 純通知, AckNone
func CastUnpinChat(msg *gonats.Msg) {
	var event protocol.UnpinChatCastEvent
	if err := json.Unmarshal(msg.Data, &event); err != nil {
		log.Printf("[CastUnpinChat] unmarshall err: %v", err)
		return
	}
	castMsg, err := json.Marshal(&protocol.Payload{MsgType: protocol.UnpinChatCast, Data: msg.Data})
	if err != nil {
		log.Printf("[CastUnpinChat] consumer marshal error: %v", err)
		return
	}
	ws.BroadcastToGroup(event.GroupID, castMsg)
}

// CastDelChat - 轉發發言刪除, 純通知, AckNone
func CastDelChat(msg *gonats.Msg) {
	var event protocol.DelChatCastEvent
	if err := json.Unmarshal(msg.Data, &event); err != nil {
		log.Printf("[CastDelChat] unmarshall err: %v", err)
		return
	}
	castMsg, err := json.Marshal(&protocol.Payload{MsgType: protocol.DelChatCast, Data: msg.Data})
	if err != nil {
		log.Printf("[CastDelChat] consumer marshal error: %v", err)
		return
	}
	ws.BroadcastToGroup(event.GroupID, castMsg)
}
