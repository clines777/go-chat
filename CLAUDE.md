# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Commands

```bash
# Build
go build ./...

# Run
go run ./cmd/chatd

# Single test
go test ./internal/ws/... -run TestFoo
```

Environment variables are loaded from `.env` at the project root. Required vars: `DB_HOST`, `DB_USER`, `DB_PASSWORD`, `DB_NAME`, `REDIS_HOST`, `REDIS_PASSWORD`.

## Architecture

Two servers run in the same process (`cmd/chatd/main.go`):
- **HTTP** `:9501` — REST API via Gin
- **WebSocket** `:9502` — real-time messaging

### WebSocket message flow

```
client → conn.ReadMessage()
       → ws.Default.Dispatch(client, input)   // routes by msg_type
       → handler returns []byte
       → client.Send <- output
       → WritePump goroutine → conn.WriteMessage()
```

Each connection gets an independent `WritePump` goroutine writing from a buffered `Send` channel, so slow clients cannot block the read loop. Error responses from the dispatcher go through `client.TrySend` (non-blocking drop if full).

WS routes are registered in `internal/handler/ws/routes.go` via `init()`, so `main.go` only needs a blank import `_ "gochat/internal/handler/ws"`.

### Session / auth model

There are three Redis token types (see `internal/protocol/redis_key.go`):
- `t_login` — short-lived, exchanged for a WS session on `login` message
- `t_api` — used by HTTP handlers for REST auth
- `t_resume` — allows reconnection without re-login (`resume` message)

WS session state is stored in Redis under `session:<connID>` as JSON (`ws.Session`: ConnID, UserID, SiteID). Routes that require auth set `SessionFree: false`; the dispatcher checks Redis before calling the handler.

### Package responsibilities

| Package | Role |
|---|---|
| `internal/ws` | Client lifecycle, hub, dispatcher, WS session lookup |
| `internal/handler/ws` | WS message handlers (one file per message type) |
| `internal/handler/api` | HTTP handlers + route registration |
| `internal/protocol` | Shared types: `Payload`, `Type` constants, Redis key helpers |
| `internal/model` | DB model structs |
| `internal/infra/db` | PostgreSQL singleton (`squirrel` query builder, `pgx` driver) |
| `internal/infra/redis` | Redis singleton |
| `internal/infra` | Env config (parsed via `caarlos0/env`) |

### Adding a new WS message type

1. Add the `Type` constant to `internal/protocol/msg_type.go`
2. Create `internal/handler/ws/<name>.go` with a `func Name(ctx *ws.Ctx) ([]byte, error)`
3. Register it in `internal/handler/ws/routes.go`
