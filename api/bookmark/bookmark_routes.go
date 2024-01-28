package bookmark

import (
	"example.com/golang_twitter/api/interfaces"
	"example.com/golang_twitter/api/middleware"
	"example.com/golang_twitter/config"
	"github.com/gin-gonic/gin"
)

func BookmarkRoutes(router *gin.RouterGroup, cfg config.Config, redisConn *interfaces.RedisConn) {
	bookmark := router.Group("/bookmarks")
	{
		bookmark.GET("", middleware.AuthRequired(), GetBookmarks(cfg, redisConn)) // ブックマーク取得機能
	}
}
