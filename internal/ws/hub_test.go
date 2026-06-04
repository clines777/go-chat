package ws

// !!讓claude幫寫的簡易handler test, 留作參考用, 之後再研究看看依賴性注入是否有額外套件可支援
// 群場景 (groupSubs) 相關函式為純記憶體操作, 不依賴任何外部基礎設施, 適合單元測試。
// 這些函式是 ban / kick / 群播功能的底層: BroadcastToGroup 群播、SendToGroupUser 定向通知、
// EvictGroupUser 將被踢者移出場景。

import "testing"

// newTestClient 建立一個可接收訊息的假連線 (Send 為 buffered channel, 不啟動 WritePump)。
func newTestClient(connID string, userID int) *Client {
	return &Client{ConnID: connID, UserId: userID, Send: make(chan []byte, 8)}
}

// recvCount 非阻塞地讀出 channel 中累積的訊息數。
func recvCount(ch chan []byte) int {
	n := 0
	for {
		select {
		case <-ch:
			n++
		default:
			return n
		}
	}
}

func TestBroadcastToGroup(t *testing.T) {
	const gid = 90001
	c1 := newTestClient("a", 1)
	c2 := newTestClient("b", 2)
	JoinGroup(c1.ConnID, gid, c1)
	JoinGroup(c2.ConnID, gid, c2)
	t.Cleanup(func() { ResetGroupScene("a"); ResetGroupScene("b") })

	BroadcastToGroup(gid, []byte("hello"))

	if got := recvCount(c1.Send); got != 1 {
		t.Errorf("c1: expected 1 message, got %d", got)
	}
	if got := recvCount(c2.Send); got != 1 {
		t.Errorf("c2: expected 1 message, got %d", got)
	}
}

func TestSendToGroupUser_OnlyTargetUser(t *testing.T) {
	const gid = 90002
	c1 := newTestClient("a", 1) // 目標 user
	c2 := newTestClient("b", 2) // 其他人
	c3 := newTestClient("c", 1) // 目標 user 的另一條連線
	JoinGroup(c1.ConnID, gid, c1)
	JoinGroup(c2.ConnID, gid, c2)
	JoinGroup(c3.ConnID, gid, c3)
	t.Cleanup(func() { ResetGroupScene("a"); ResetGroupScene("b"); ResetGroupScene("c") })

	SendToGroupUser(gid, 1, []byte("ban"))

	if got := recvCount(c1.Send); got != 1 {
		t.Errorf("c1 (target): expected 1, got %d", got)
	}
	if got := recvCount(c2.Send); got != 0 {
		t.Errorf("c2 (other user): expected 0, got %d", got)
	}
	if got := recvCount(c3.Send); got != 1 {
		t.Errorf("c3 (target's 2nd conn): expected 1, got %d", got)
	}
}

func TestEvictGroupUser_RemovesFromBroadcast(t *testing.T) {
	const gid = 90003
	c1 := newTestClient("a", 1) // 被踢者
	c2 := newTestClient("b", 2) // 留下的成員
	c3 := newTestClient("c", 1) // 被踢者的另一條連線
	JoinGroup(c1.ConnID, gid, c1)
	JoinGroup(c2.ConnID, gid, c2)
	JoinGroup(c3.ConnID, gid, c3)
	t.Cleanup(func() { ResetGroupScene("a"); ResetGroupScene("b"); ResetGroupScene("c") })

	conns := EvictGroupUser(gid, 1)
	if len(conns) != 2 {
		t.Fatalf("expected 2 evicted conns (a, c), got %d: %v", len(conns), conns)
	}
	gotSet := map[string]bool{}
	for _, id := range conns {
		gotSet[id] = true
	}
	if !gotSet["a"] || !gotSet["c"] {
		t.Errorf("expected evicted conns to be {a, c}, got %v", conns)
	}

	// 踢出後群播只會送達留下的成員
	BroadcastToGroup(gid, []byte("after-kick"))
	if got := recvCount(c1.Send); got != 0 {
		t.Errorf("evicted c1: expected 0, got %d", got)
	}
	if got := recvCount(c2.Send); got != 1 {
		t.Errorf("remaining c2: expected 1, got %d", got)
	}
	if got := recvCount(c3.Send); got != 0 {
		t.Errorf("evicted c3: expected 0, got %d", got)
	}
}

func TestEvictGroupUser_NoMatch(t *testing.T) {
	const gid = 90004
	c1 := newTestClient("a", 1)
	JoinGroup(c1.ConnID, gid, c1)
	t.Cleanup(func() { ResetGroupScene("a") })

	if conns := EvictGroupUser(gid, 999); len(conns) != 0 {
		t.Errorf("expected no eviction for absent user, got %v", conns)
	}
	// 原連線仍應正常收到群播
	BroadcastToGroup(gid, []byte("x"))
	if got := recvCount(c1.Send); got != 1 {
		t.Errorf("c1 should still be in group: expected 1, got %d", got)
	}
}

func TestResetGroupScene(t *testing.T) {
	const gid = 90005
	c1 := newTestClient("a", 1)
	JoinGroup(c1.ConnID, gid, c1)

	ResetGroupScene(c1.ConnID)
	BroadcastToGroup(gid, []byte("x"))
	if got := recvCount(c1.Send); got != 0 {
		t.Errorf("after reset c1 should receive nothing, got %d", got)
	}
}

// TestJoinGroup_SwitchRemovesOld 同一連線切換群組後, 舊群的群播不應再送達。
func TestJoinGroup_SwitchRemovesOld(t *testing.T) {
	const oldGid, newGid = 90006, 90007
	c := newTestClient("a", 1)
	JoinGroup(c.ConnID, oldGid, c)
	JoinGroup(c.ConnID, newGid, c) // 切換場景
	t.Cleanup(func() { ResetGroupScene("a") })

	BroadcastToGroup(oldGid, []byte("old"))
	if got := recvCount(c.Send); got != 0 {
		t.Errorf("old group should no longer reach the conn, got %d", got)
	}
	BroadcastToGroup(newGid, []byte("new"))
	if got := recvCount(c.Send); got != 1 {
		t.Errorf("new group should reach the conn, got %d", got)
	}
}
