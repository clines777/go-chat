package consumer

import (
	"encoding/json"
	gonats "github.com/nats-io/nats.go"
	"gochat/internal/protocol"
	"gochat/internal/session"
	"gochat/internal/ws"
	"log"
)

// CastLeaveGroup 轉發用戶退群通知. 純通知, AckNone
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

// CastDelChat - 轉發聊天訊息刪除, 純通知, AckNone
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

// CastBanUser - 通知被禁言的 user (其連線目前在該群), 純通知, AckNone
func CastBanUser(msg *gonats.Msg) {
	var event protocol.BanUserEvent
	if err := json.Unmarshal(msg.Data, &event); err != nil {
		log.Printf("[CastBanUser] unmarshall err: %v", err)
		return
	}
	ackMsg, err := json.Marshal(&protocol.Payload{MsgType: protocol.BanUserAck, Data: msg.Data})
	if err != nil {
		log.Printf("[CastBanUser] consumer marshal error: %v", err)
		return
	}
	ws.SendToGroupUser(event.GroupID, event.UserID, ackMsg)
}

// CastUnbanUser - 通知被解除禁言的 user, 純通知, AckNone
func CastUnbanUser(msg *gonats.Msg) {
	var event protocol.UnbanUserEvent
	if err := json.Unmarshal(msg.Data, &event); err != nil {
		log.Printf("[CastUnbanUser] unmarshall err: %v", err)
		return
	}
	ackMsg, err := json.Marshal(&protocol.Payload{MsgType: protocol.UnbanUserAck, Data: msg.Data})
	if err != nil {
		log.Printf("[CastUnbanUser] consumer marshal error: %v", err)
		return
	}
	ws.SendToGroupUser(event.GroupID, event.UserID, ackMsg)
}

// CastKickUser - 通知被踢出的 user 並將其移出群場景, 再向群內其餘成員廣播, 純通知, AckNone
func CastKickUser(msg *gonats.Msg) {
	var event protocol.KickUserEvent
	if err := json.Unmarshal(msg.Data, &event); err != nil {
		log.Printf("[CastKickUser] unmarshall err: %v", err)
		return
	}

	ackMsg, err := json.Marshal(&protocol.Payload{MsgType: protocol.KickUserAck, Data: msg.Data})
	if err != nil {
		log.Printf("[CastKickUser] consumer marshal error: %v", err)
		return
	}
	// 先通知被踢者, 再將其連線移出群場景並重設 session
	ws.SendToGroupUser(event.GroupID, event.UserID, ackMsg)
	for _, connID := range ws.EvictGroupUser(event.GroupID, event.UserID) {
		sess := session.Get(event.UserID, connID)
		if sess == nil {
			continue
		}
		sess.InGroupID = 0
		sess.Scene = protocol.SceneMyGroup
		if err := session.Set(event.UserID, connID, sess); err != nil {
			log.Printf("[CastKickUser] session.Set user=%d conn=%s error: %v", event.UserID, connID, err)
		}
	}

	castMsg, err := json.Marshal(&protocol.Payload{MsgType: protocol.KickUserCast, Data: msg.Data})
	if err != nil {
		log.Printf("[CastKickUser] cast marshal error: %v", err)
		return
	}
	ws.BroadcastToGroup(event.GroupID, castMsg)
}
