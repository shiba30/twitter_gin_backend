package notification

import (
	"log"

	"example.com/golang_twitter/api/interfaces"
	"example.com/golang_twitter/config"
	sqlc "example.com/golang_twitter/db/sqlc"
	"example.com/golang_twitter/utils"
	"github.com/gin-gonic/gin"
)

func NotificationRoutes(router *gin.RouterGroup, cfg config.Config, redisConn *interfaces.RedisConn, queries *sqlc.Queries) {
	router.GET("/notifications", GetNotifications(cfg, redisConn, queries))
}

func GetNotifications(cfg config.Config, redisConn *interfaces.RedisConn, queries *sqlc.Queries) gin.HandlerFunc {
	return func(c *gin.Context) {
		// ユーザID取得
		userInfo, err := utils.CurrentUser(c, redisConn, queries)
		if err != nil {
			log.Printf("Failed to get userInfo: %v", err)
			c.JSON(500, gin.H{"error": "プロフィール情報の取得に失敗しました"})
			return
		}

		// 通知情報を取得
		notifications, err := queries.GetNotifications(c, userInfo.ID)
		if err != nil {
			log.Printf("Failed to get notifications: %v", err)
			c.JSON(500, gin.H{"error": "通知情報の取得に失敗しました"})
			return
		}

		/*
			// 通知情報を既読にする
			_, err = queries.UpdateNotifications(c, userInfo.ID)
			if err != nil {
				log.Printf("Failed to update notifications: %v", err)
				c.JSON(500, gin.H{"error": "通知情報の更新に失敗しました"})
				return
			}
		*/

		log.Printf("Get notifications: %v", notifications)

		c.JSON(200, gin.H{
			"notifications": notifications,
		})
	}
}
