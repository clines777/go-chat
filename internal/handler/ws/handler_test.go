package handler

// !!讓claude幫寫的簡易handler test, 留作參考用, 之後再研究看看依賴性注入是否有額外套件可支援
// 這些測試只覆蓋 handler 在「呼叫任何基礎設施 (Redis / DB / NATS) 之前」就會返回的分支,
// 也就是輸入驗證 (msg payload 解析與必填欄位檢查)。這些路徑不需要連線, 可純單元測試。
//
// 需要 session / DB / NATS 的成功路徑屬於整合測試範疇 (handler 目前直接依賴全域單例,
// 沒有注入點), 不在此檔涵蓋。

import (
	"encoding/json"
	"testing"

	"gochat/internal/protocol"
	"gochat/internal/ws"
)

// newCtx 建立一個帶有指定 payload data 的 Ctx。data 為空字串時 Data 留空 (模擬缺少 body)。
func newCtx(msgType protocol.MsgType, data string) *ws.Ctx {
	var raw json.RawMessage
	if data != "" {
		raw = json.RawMessage(data)
	}
	return &ws.Ctx{
		Client:  &ws.Client{UserId: 1, ConnID: "test-conn"},
		Payload: &protocol.Payload{MsgType: msgType, Data: raw},
	}
}

// assertInvalidParam 驗證 handler 回傳的是 ErrInvalidParam 錯誤 payload。
func assertInvalidParam(t *testing.T, got *protocol.Payload) {
	t.Helper()
	if got == nil {
		t.Fatal("expected an error payload, got nil")
	}
	if got.MsgType != protocol.Error {
		t.Fatalf("expected MsgType=%q, got %q", protocol.Error, got.MsgType)
	}
	if got.ErrCode != protocol.ErrInvalidParam {
		t.Fatalf("expected ErrCode=%d (ErrInvalidParam), got %d", protocol.ErrInvalidParam, got.ErrCode)
	}
	if got.Origin == nil {
		t.Error("expected Origin to reference the original payload, got nil")
	}
}

// TestHandlers_RejectInvalidInput 對每個 handler 測試: 缺少 body、格式錯誤、缺必填欄位 → ErrInvalidParam。
func TestHandlers_RejectInvalidInput(t *testing.T) {
	cases := []struct {
		name    string
		handler func(*ws.Ctx) *protocol.Payload
		msgType protocol.MsgType
		data    string
	}{
		// ── BanUser (需要 group_id + user_id) ──
		{"BanUser/missingBody", BanUser, protocol.BanUser, ""},
		{"BanUser/malformed", BanUser, protocol.BanUser, "{"},
		{"BanUser/empty", BanUser, protocol.BanUser, `{}`},
		{"BanUser/missingUser", BanUser, protocol.BanUser, `{"group_id":5}`},
		{"BanUser/missingGroup", BanUser, protocol.BanUser, `{"user_id":7}`},

		// ── UnbanUser (需要 group_id + user_id) ──
		{"UnbanUser/missingBody", UnbanUser, protocol.UnbanUser, ""},
		{"UnbanUser/malformed", UnbanUser, protocol.UnbanUser, "}"},
		{"UnbanUser/missingUser", UnbanUser, protocol.UnbanUser, `{"group_id":5}`},
		{"UnbanUser/missingGroup", UnbanUser, protocol.UnbanUser, `{"user_id":7}`},

		// ── KickUser (需要 group_id + user_id) ──
		{"KickUser/missingBody", KickUser, protocol.KickUser, ""},
		{"KickUser/malformed", KickUser, protocol.KickUser, `{"group_id":}`},
		{"KickUser/missingUser", KickUser, protocol.KickUser, `{"group_id":5}`},
		{"KickUser/missingGroup", KickUser, protocol.KickUser, `{"user_id":7}`},

		// ── PinChat (需要 group_id + chat_id) ──
		{"PinChat/missingBody", PinChat, protocol.PinChat, ""},
		{"PinChat/empty", PinChat, protocol.PinChat, `{}`},
		{"PinChat/missingChat", PinChat, protocol.PinChat, `{"group_id":5}`},
		{"PinChat/missingGroup", PinChat, protocol.PinChat, `{"chat_id":9}`},

		// ── UnpinChat (需要 group_id) ──
		{"UnpinChat/missingBody", UnpinChat, protocol.UnpinChat, ""},
		{"UnpinChat/empty", UnpinChat, protocol.UnpinChat, `{}`},
		{"UnpinChat/malformed", UnpinChat, protocol.UnpinChat, "nope"},

		// ── DelChat (需要 group_id + chat_id) ──
		{"DelChat/missingBody", DelChat, protocol.DelChat, ""},
		{"DelChat/missingChat", DelChat, protocol.DelChat, `{"group_id":5}`},
		{"DelChat/missingGroup", DelChat, protocol.DelChat, `{"chat_id":9}`},

		// ── EnterGroup (需要 group_id) ──
		{"EnterGroup/missingBody", EnterGroup, protocol.EnterGroup, ""},
		{"EnterGroup/empty", EnterGroup, protocol.EnterGroup, `{}`},

		// ── GetHistory (需要 group_id; before_id 選填) ──
		{"GetHistory/missingBody", GetHistory, protocol.GetHistory, ""},
		{"GetHistory/empty", GetHistory, protocol.GetHistory, `{}`},
		{"GetHistory/malformed", GetHistory, protocol.GetHistory, `{"group_id":}`},
		{"GetHistory/zeroGroup", GetHistory, protocol.GetHistory, `{"group_id":0,"before_id":100}`},

		// ── JoinGroup (需要 group_id > 0) ──
		{"JoinGroup/missingBody", JoinGroup, protocol.JoinGroup, ""},
		{"JoinGroup/zero", JoinGroup, protocol.JoinGroup, `{"group_id":0}`},
		{"JoinGroup/negative", JoinGroup, protocol.JoinGroup, `{"group_id":-3}`},

		// ── LeaveGroup (需要 group_id > 0) ──
		{"LeaveGroup/missingBody", LeaveGroup, protocol.LeaveGroup, ""},
		{"LeaveGroup/zero", LeaveGroup, protocol.LeaveGroup, `{"group_id":0}`},

		// ── SendChat (需要 group_id + content) ──
		{"SendChat/missingBody", SendChat, protocol.SendChat, ""},
		{"SendChat/missingContent", SendChat, protocol.SendChat, `{"group_id":5}`},
		{"SendChat/missingGroup", SendChat, protocol.SendChat, `{"content":"hi"}`},
		{"SendChat/emptyContent", SendChat, protocol.SendChat, `{"group_id":5,"content":""}`},

		// ── UpdateGroup (需要 group_id, 驗證失敗路徑) ──
		{"UpdateGroup/missingBody", UpdateGroup, protocol.UpdateGroup, ""},
		{"UpdateGroup/missingGroup", UpdateGroup, protocol.UpdateGroup, `{"title":"x"}`},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			got := c.handler(newCtx(c.msgType, c.data))
			assertInvalidParam(t, got)
		})
	}
}

// TestUpdateGroup_NoFields 當 group_id 合法但 title/bulletin/remark 全空時, 視為無變更直接回 UpdateGroupOk,
// 此分支在呼叫 session / DB 之前返回, 可純單元測試。
func TestUpdateGroup_NoFields(t *testing.T) {
	got := UpdateGroup(newCtx(protocol.UpdateGroup, `{"group_id":5}`))
	if got == nil {
		t.Fatal("expected a payload, got nil")
	}
	if got.MsgType != protocol.UpdateGroupOk {
		t.Fatalf("expected MsgType=%q, got %q (err_code=%d)", protocol.UpdateGroupOk, got.MsgType, got.ErrCode)
	}
}

// TestFindPinChat_NoPin pinChatID <= 0 時直接回 nil, 不查 DB。
func TestFindPinChat_NoPin(t *testing.T) {
	if got := FindPinChat(10, 0); got != nil {
		t.Errorf("pinChatID=0: expected nil, got %+v", got)
	}
	if got := FindPinChat(10, -1); got != nil {
		t.Errorf("pinChatID<0: expected nil, got %+v", got)
	}
}
