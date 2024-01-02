package reply

import (
	"example.com/golang_twitter/api/middleware"
	"example.com/golang_twitter/config"
	"github.com/gin-gonic/gin"
)

func ReplyRoutes(router *gin.RouterGroup, cfg config.Config) {
	tweet := router.Group("/reply")
	{
		tweet.POST("/:id", middleware.AuthRequired(), postReply(cfg))       // コメント投稿機能
		tweet.GET("/replies/:id", middleware.AuthRequired(), getReply(cfg)) // コメント一覧取得機能
	}
}
