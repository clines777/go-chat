package chat

// trimHistoryPage 是向上分頁「多取一筆」策略的純邏輯部分 (不碰 DB), 可單元測試。
// GetHistoryChats 以 LIMIT limit+1 查詢, 由此函式判斷是否還有更舊訊息並修剪回 limit 筆。

import (
	"testing"

	"gochat/internal/protocol"
)

// ids 取出 ChatInfo 的 id 序列, 方便比對。
func ids(chats []protocol.ChatInfo) []int {
	out := make([]int, len(chats))
	for i, c := range chats {
		out[i] = c.ID
	}
	return out
}

// page 產生 id 正序 (1..n) 的測試資料, 模擬 DB 回傳。
func page(n int) []protocol.ChatInfo {
	rows := make([]protocol.ChatInfo, n)
	for i := 0; i < n; i++ {
		rows[i] = protocol.ChatInfo{ID: i + 1}
	}
	return rows
}

func TestTrimHistoryPage(t *testing.T) {
	const limit = 30

	cases := []struct {
		name        string
		rows        []protocol.ChatInfo
		limit       int
		wantHasMore bool
		wantIDs     []int
	}{
		{
			name:        "empty/沒有任何訊息",
			rows:        page(0),
			limit:       limit,
			wantHasMore: false,
			wantIDs:     []int{},
		},
		{
			name:        "fewerThanLimit/不足一頁",
			rows:        page(5),
			limit:       limit,
			wantHasMore: false,
			wantIDs:     []int{1, 2, 3, 4, 5},
		},
		{
			name:        "exactlyLimit/剛好一頁不算還有更多",
			rows:        page(limit),
			limit:       limit,
			wantHasMore: false,
			wantIDs:     seq(1, limit),
		},
		{
			name:        "oneExtra/多取一筆代表還有更舊, 砍掉最舊那筆",
			rows:        page(limit + 1),
			limit:       limit,
			wantHasMore: true,
			wantIDs:     seq(2, limit+1), // id=1 (最舊) 被砍, 保留較新的 limit 筆
		},
		{
			name:        "limitOne/邊界 limit=1 且有更多",
			rows:        page(2),
			limit:       1,
			wantHasMore: true,
			wantIDs:     []int{2},
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			got, hasMore := trimHistoryPage(c.rows, c.limit)
			if hasMore != c.wantHasMore {
				t.Errorf("hasMore = %v, want %v", hasMore, c.wantHasMore)
			}
			if gotIDs := ids(got); !equalInts(gotIDs, c.wantIDs) {
				t.Errorf("ids = %v, want %v", gotIDs, c.wantIDs)
			}
		})
	}
}

func seq(lo, hi int) []int {
	out := make([]int, 0, hi-lo+1)
	for i := lo; i <= hi; i++ {
		out = append(out, i)
	}
	return out
}

func equalInts(a, b []int) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}
