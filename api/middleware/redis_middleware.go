package middleware

import (
	"example.com/golang_twitter/api/interfaces"
	"github.com/gin-gonic/gin"
)

// redis接続をコンテキストに追加する関数
func RedisMiddleware(redisConn *interfaces.RedisConn) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("redisConn", redisConn)
		c.Next()
	}
}
