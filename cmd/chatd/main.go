package main

import (
	"context"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"gochat/internal/handler/api"
	_ "gochat/internal/handler/ws"
	"gochat/internal/infra/db"
	"gochat/internal/infra/redis"
	"gochat/internal/middleware"
	"gochat/internal/ws"
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

	client := ws.NewClient(conn)
	go client.WritePump()

	conn.SetPongHandler(func(string) error {
		_ = conn.SetReadDeadline(time.Now().Add(60 * time.Second))
		return nil
	})

	for {
		mt, input, err := conn.ReadMessage()
		if err != nil {
			log.Printf("[WS] read error: %v", err)
			close(client.Send)
			return
		}
		if mt != websocket.TextMessage {
			ws.Unregister(client.ConnID)
			close(client.Send)
			client.Conn.Close()
			return
		}

		output, dispErr := ws.Default.Dispatch(client, input)
		if dispErr != nil {
			ws.Unregister(client.ConnID)
			client.Conn.Close()
			close(client.Send)
			return
		}

		if len(output) > 0 {
			client.Send <- output
		}
	}
}
