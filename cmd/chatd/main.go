package main

import (
	"context"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"gochat/internal/handler/api"
	"gochat/internal/infra/db"
	"gochat/internal/infra/redis"
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

	//BEGIN 初始化DB連線池
	dbConn, dbErr := db.GetDBConn()
	if dbErr != nil {
		log.Fatalf("db init failed: %v", dbErr)
	}

	if err := dbConn.Ping(); err != nil {
		log.Fatalf("db health check failed: %v", err)
	}
	log.Println("[DB] check success")
	//END 初始化DB連線池

	//BEGIN 初始化Redis連線池
	redisConn := redis.GetRedis()

	redisPingErr := redisConn.Ping()
	if redisPingErr != nil {
		log.Fatalf("redisConn health check failed: %v", redisPingErr)
	}
	//END 初始化Redis連線池

	//BEGIN 初始化Http server
	r := gin.New()

	api.RegisterRoutes(r)

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
	//END 初始化Http server

	//BEGIN 初始化WS server
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
	//END 初始化WS server

	//BEGIN 註冊server down事件
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("shutting down...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_ = httpSrv.Shutdown(ctx)
	_ = wsSrv.Shutdown(ctx)
	log.Println("server closed")
	//END 註冊server down行為
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
		mt, input, err := conn.ReadMessage()
		if err != nil {
			log.Printf("[WS] read error: %v", err)
			return
		}

		if mt != websocket.TextMessage {
			conn.Close()
			return
		}

		dispatcher := ws.NewDispatcher()

		output, err := dispatcher.Dispatch(input)

		if err := conn.WriteMessage(mt, output); err != nil {
			log.Printf("[WS] write error: %v", err)
			return
		}
	}
}
