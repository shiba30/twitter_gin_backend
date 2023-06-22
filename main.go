package main

import (
	"github.com/gin-gonic/gin"
)

func main() {

	router := gin.Default()

	router.GET("/health_check", func(c *gin.Context) {
		c.IndentedJSON(200, gin.H{"status": "ok"})
	})
	router.Run("0.0.0.0:8080")
}
