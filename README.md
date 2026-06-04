# gochat

基於 WebSocket 的即時群組聊天服務。單一行程同時跑兩個 server:HTTP (`:9501`) 提供 REST API,WebSocket (`:9502`) 處理即時訊息;群組事件透過 NATS JetStream 廣播、session 存於 Redis。
由於開發目的僅用於demo，架構直接採all in one.

## 環境需求

 - DB: PostgreSQL >= 8.2
 - Redis >= 7.2
 - nats >= 2.14 (啟用jet-stream)
 - Go >= 1.25

## 設定

設定由專案根目錄的 `.env` 載入,主要變數:

| 變數 | 說明                  |
|---|---------------------|
| `DB_HOST` / `DB_PORT` / `DB_NAME` / `DB_USER` / `DB_PASSWORD` | PostgreSQL 連線       |
| `DB_SSLMODE` / `DB_TIMEZONE` | 連線模式與時區             |
| `DB_MAX_OPEN` / `DB_MAX_IDLE` / `DB_CONN_MAX_LIFETIME` / `DB_CONN_MAX_IDLE` | 連線池設定               |
| `REDIS_HOST` / `REDIS_PASSWORD` / `REDIS_DB` / `REDIS_DIAL_TIMEOUT` | Redis 連線            |
| `NATS_URL` | NATS 連線位址           |
| `SERVER_NAME` | 此節點名稱 (多台server識別用) |

## 初始化 DB

首次啟動前先建表 (schema 見 `db/schema.sql`):

```bash
psql "$DB_URL" -f db/schema.sql
```

## 啟動

```bash
go run ./cmd/chatd/main.go
```

##前端頁面: 全程由claude編寫, 
 - 訪問位置:localhost:9501/
 - 毋需驗證, 隨意填入username即可登入
