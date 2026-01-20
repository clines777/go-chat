package main

import (
	"context"
	"errors"
	"github.com/caarlos0/env/v11"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/joho/godotenv"
	"gochat/internal/config"
	"gochat/internal/db"
	"gochat/internal/redis"
	"gochat/internal/route"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true },
}

func main() {

	// 本地開發環境檢查並載入env參數到os
	envPath := "../.env"
	_, envFileErr := os.Stat(envPath)
	if !os.IsNotExist(envFileErr) {
		_ = godotenv.Load()
	}

	//載入env
	var cfg config.EnvConfig
	if cfgErr := env.Parse(&cfg); cfgErr != nil {
		log.Fatalf("parse env failed: %v", cfgErr)
	}

	dbConn, dbErr := db.Init(&cfg)
	if dbErr != nil {
		log.Fatalf("db init failed: %v", dbErr)
	}

	defer dbConn.Close()

	pingCtx, pingCancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer pingCancel()

	if err := dbConn.Ping(pingCtx); err != nil {
		log.Fatalf("db health check failed: %v", err)
	}
	log.Println("[DB] check success")

	redisConn, redisErr := redis.Init(&cfg)
	if redisErr != nil {
		log.Fatalf("redisConn init failed: %v", redisErr)
	}

	defer redisConn.Close()

	redisPingCtx, redisPingCancel := context.WithTimeout(context.Background(), 3*time.Second)
	redisConn.Ping(redisPingCtx)

	defer redisPingCancel()

	//on Http server
	r := gin.New()

	route.RegisterRoutes(r)

	r.Use(gin.Logger(), gin.Recovery())

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

	//on Websocket server
	wsMux := http.NewServeMux()
	wsMux.HandleFunc("/ws", wsHandler)

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

	//server close action
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

func wsHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("[WS] upgrade error: %v", err)
		return
	}
	defer conn.Close()

	conn.SetPongHandler(func(string) error {
		_ = conn.SetReadDeadline(time.Now().Add(60 * time.Second))
		return nil
	})
	_ = conn.SetReadDeadline(time.Now().Add(60 * time.Second))

	for {
		mt, msg, err := conn.ReadMessage()
		if err != nil {
			log.Printf("[WS] read error: %v", err)
			return
		}

		if err := conn.WriteMessage(mt, msg); err != nil {
			log.Printf("[WS] write error: %v", err)
			return
		}
	}
}
