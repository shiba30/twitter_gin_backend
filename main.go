package main

import (
	"log"

	"example.com/golang_twitter/api"
	"example.com/golang_twitter/db"
	"github.com/gin-gonic/gin"
)

// Twitter clone by golang(gin)
func main() {

	// DB接続
	db.ConnectDB()
	defer db.DbConn().Close()

	router := gin.Default()

	router.Static("/static", "./static")
	router.LoadHTMLGlob("templates/*")

	router.GET("/health_check", func(c *gin.Context) {
		c.IndentedJSON(200, gin.H{"status": "ok"})
	})

	api.Routes(router)

	log.Println("Server started!")
	router.Run("0.0.0.0:8080")
}
