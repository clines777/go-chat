# syntax=docker/dockerfile:1

# ---------- build stage ----------
FROM golang:1.25-alpine AS builder
WORKDIR /src

# 先抓相依，利用 layer cache
COPY go.mod go.sum ./
RUN go mod download

COPY . .
# 靜態連結，產出單一 binary
RUN CGO_ENABLED=0 GOOS=linux go build -trimpath -ldflags="-s -w" -o /out/chatd ./cmd/chatd

# ---------- runtime stage ----------
FROM alpine:3.20
RUN apk add --no-cache ca-certificates tzdata wget
WORKDIR /app

COPY --from=builder /out/chatd /app/chatd
# app 從工作目錄讀 ./web/index.html 與 ./static
COPY web ./web
COPY static ./static

EXPOSE 9501 9502

# health check: HTTP server 的 "/" 會回 index.html
HEALTHCHECK --interval=15s --timeout=3s --retries=5 \
  CMD wget --spider -q http://127.0.0.1:9501/ || exit 1

ENTRYPOINT ["/app/chatd"]
