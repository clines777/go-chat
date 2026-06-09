# gochat

基於 WebSocket 的即時群聊服務。參考個人之前開發的一套基於PHP(Hyperf框架)、Redis、Mysql，以包網平台用戶為目標的群聊系統為基礎，採用Go、Nats、Redis與Postgresql開發而成，因應demo需求，對部分機制加以簡化，並將架構調整為可直接在本機運行的單體式架構。
此專案開發主要作為面試作品用，同時也是個人為練習之前沒嘗試過的AI coding的開發方式而做，後端部分為手動編程與Claude CLI交互開發而成，claude code主要用於code review、協助進行大規模調整及提供優化建議，附帶的web前端頁面則全由AI處理。

## 環境需求

 - DB: Postgres >= 8.2
 - Redis >= 7.2
 - nats >= 2.14 (啟用JetStream)
 - Go >= 1.25

## 技術選型因素

- Postgres
1. 查詢優化器可預測，不容易出現rule-based跟cost-based錯用的情況。
2. 已有可靠的extension - Citus，可用於日後DB數據量過大時方便做sharding。
3. 豐富的extension庫。
- Redis
1. 聲帶成熟穩定泛用的no sql。
2. 豐富的數據結構方便存儲用戶狀態。
3. TTL方便用於短生命週期的數據keep alive，例如用戶session跟api token。
- Nats(JetStream)
1. 作為拉式MQ，方便控流量。
2. 規模較小，安裝部屬方便。
3. subject(類似Kafka的Topic)匹配的wild card方便做數據分割。
4. 自帶JetStream streaming layer模組，相較於傳統的pub/sub，(1)應用性較廣 (2)debug時方便觀察管道內狀況。(3)資料持久化

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

## 前端頁面: 
 - 訪問位置:localhost:9501/
 - 毋需驗證, 隨意填入username即可登入
