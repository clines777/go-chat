package main

import (
	"context"
	"errors"
	"gochat/internal/ws"

	"github.com/gin-gonic/gin"
	"gochat/internal/chat"
	"gochat/internal/handler/api"
	_ "gochat/internal/handler/ws"
	"gochat/internal/infra"
	"gochat/internal/infra/db"
	infranats "gochat/internal/infra/nats"
	"gochat/internal/infra/redis"
	"gochat/internal/middleware"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {

	dbConn, dbErr := db.GetDBConn()
	if dbErr != nil {
		log.Fatalf("db init failed: %v", dbErr)
	}
	if err := dbConn.Ping(); err != nil {
		log.Fatalf("db health check failed: %v", err)
	}
	log.Println("[DB] check success")

	redisConn := redis.GetRedis()
	if err := redisConn.Ping(); err != nil {
		log.Fatalf("redis health check failed: %v", err)
	}
	log.Println("[REDIS] check success")

	natsClient, natsErr := infranats.GetNats()
	if natsErr != nil {
		log.Fatalf("nats init failed: %v", natsErr)
	}
	if err := natsClient.Ping(); err != nil {
		log.Fatalf("nats health check failed: %v", err)
	}
	if err := natsClient.EnsureStream(); err != nil {
		log.Fatalf("nats stream setup failed: %v", err)
	}
	log.Println("[NATS] check success")

	serverName := infra.GetEnvConfig().ServerName
	//subscribe群組聊天管道
	if err := chat.StartGroupChatConsumer(serverName); err != nil {
		log.Fatalf("nats consumer start failed: %v", err)
	}

	r := gin.New()
	r.Use(middleware.Logger(), gin.Recovery())
	api.RegisterRoutes(r)

	httpSrv := &http.Server{
		Addr:    ":9501",
		Handler: r,
	}
	go func() {
		log.Println("[HTTP] listening on :9501")
		if err := httpSrv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("[HTTP] listen error: %v", err)
		}
	}()

	wsMux := http.NewServeMux()
	wsMux.HandleFunc("/ws", ws.HandleWs)
	wsSrv := &http.Server{
		Addr:    ":9502",
		Handler: wsMux,
	}
	go func() {
		log.Println("[WS] listening on :9502 (path: /ws)")
		if err := wsSrv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("[WS] listen error: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("shutting down...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_ = httpSrv.Shutdown(ctx)
	_ = wsSrv.Shutdown(ctx)
	log.Println("server closed")
}
