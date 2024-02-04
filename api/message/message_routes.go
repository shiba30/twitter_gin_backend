package message

import (
	"example.com/golang_twitter/api/interfaces"
	"example.com/golang_twitter/config"
	sqlc "example.com/golang_twitter/db/sqlc"
	"github.com/gin-gonic/gin"
)

func MessageRoutes(router *gin.RouterGroup, cfg config.Config, redisConn *interfaces.RedisConn, queries *sqlc.Queries) {
	message := router.Group("/message")
	{
		message.GET("", GetMessagePage(cfg, redisConn, queries))        // メッセージページの取得
		message.POST("", PostMessage(cfg, redisConn, queries))          // メッセージの送信
		message.GET("/ws", UpgradeToWebSocket(cfg, redisConn, queries)) // WebSocketへのアップグレード
	}
}
