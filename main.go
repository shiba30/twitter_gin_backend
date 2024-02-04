package main

import (
	"log"

	"example.com/golang_twitter/api"
	"example.com/golang_twitter/api/interfaces"
	"example.com/golang_twitter/api/message"
	"example.com/golang_twitter/api/middleware"
	"example.com/golang_twitter/config"
	"example.com/golang_twitter/db"
	sqlc "example.com/golang_twitter/db/sqlc"
	"github.com/gin-gonic/gin"
)

// Twitter clone by golang(gin)
func main() {
	// 環境変数取得
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	// DB接続
	if err := db.ConnectDB(cfg); err != nil {
		log.Fatalf("failed to connect to DB: %v", err)
	}
	defer db.DbConn().Close()

	// Redis接続を初期化
	redisConn := interfaces.NewRedisConn(cfg)
	if err != nil {
		log.Fatalf("failed to initialize Redis: %v", err)
	}
	defer redisConn.Close()

	// SQLCクエリオブジェクトの生成
	queries := sqlc.New(db.DbConn())

	// WebSocketメッセージハンドリングゴルーチンの起動
	go message.HandleMessages(queries)

	router := gin.Default()

	// redisConnを全てのルートで使用
	router.Use(middleware.RedisMiddleware(redisConn))

	router.Static("/static", "./static")
	router.LoadHTMLGlob("templates/*")

	router.GET("/health_check", func(c *gin.Context) {
		c.IndentedJSON(200, gin.H{"status": "ok"})
	})

	controller := api.NewController(cfg, redisConn, queries)
	controller.Routes(router)

	log.Println("Server started!")
	router.Run(cfg.ServerAddress)
}
