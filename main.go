package main

import (
	"log"

	"example.com/golang_twitter/api"
	"example.com/golang_twitter/config"
	"example.com/golang_twitter/db"
	"github.com/gin-gonic/gin"
)

// Twitter clone by golang(gin)
func main() {
	// 環境変数取得
	cfg := config.LoadConfig()

	// DB接続
	db.ConnectDB(cfg)
	defer db.DbConn().Close()

	router := gin.Default()

	router.Static("/static", "./static")
	router.LoadHTMLGlob("templates/*")

	router.GET("/health_check", func(c *gin.Context) {
		c.IndentedJSON(200, gin.H{"status": "ok"})
	})

	controller := api.NewController(cfg)
	controller.Routes(router)

	log.Println("Server started!")
	router.Run("0.0.0.0:8080")
}
