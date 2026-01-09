package chatd

import (
	"context"
	"errors"
	"github.com/caarlos0/env/v11"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/joho/godotenv"
	"gochat/internal/db"
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

	//DB init
	_ = godotenv.Load()
	var cfg db.Config
	if err := env.Parse(&cfg); err != nil {
		log.Fatalf("parse env failed: %v", err)
	}

	d, err := db.New(db.Config{
		Host:            cfg.Host,
		Port:            cfg.Port,
		User:            cfg.User,
		Password:        cfg.Password,
		DBName:          cfg.DBName,
		SSLMode:         cfg.SSLMode,
		TimeZone:        cfg.TimeZone,
		MaxOpenConn:     cfg.MaxOpenConn,
		MaxIdleConn:     cfg.MaxIdleConn,
		ConnMaxLifetime: cfg.ConnMaxLifetime,
		ConnMaxIdleTime: cfg.ConnMaxIdleTime,
	})
	if err != nil {
		log.Fatalf("db init failed: %v", err)
	}

	defer d.Close()

	// DB health check with timeout (fail fast)
	pingCtx, pingCancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer pingCancel()

	if err := d.Ping(pingCtx); err != nil {
		log.Fatalf("db health check failed: %v", err)
	}
	log.Println("[DB] connected and healthy.")

	//on Http server
	r := gin.New()
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

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

	// 簡單 echo + ping/pong 範例
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
		// echo 回去
		if err := conn.WriteMessage(mt, msg); err != nil {
			log.Printf("[WS] write error: %v", err)
			return
		}
	}
}
