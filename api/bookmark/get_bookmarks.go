package bookmark

import (
	"context"
	"log"
	"net/http"

	"example.com/golang_twitter/api/interfaces"
	"example.com/golang_twitter/config"
	sqlc "example.com/golang_twitter/db/sqlc"
	"example.com/golang_twitter/utils"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
)

type RedisGetter interface {
	Get(ctx context.Context, key string) *redis.StringCmd
}

func GetBookmarks(cfg config.Config, redisConn *interfaces.RedisConn, queries *sqlc.Queries) gin.HandlerFunc {
	return func(c *gin.Context) {
		// セッションからユーザー情報を取得
		sessionID, err := c.Cookie("session_id")
		if err != nil {
			log.Printf("Failed to retrieve session ID: %v", err)
			c.JSON(500, gin.H{"error": "セッションIDの取得に失敗しました"})
			return
		}
		userID, err := utils.GetSessionUserId(c, sessionID)
		if err != nil {
			log.Printf("Failed to retrieve userId from session: %v", err)
			c.JSON(500, gin.H{"error": "セッションからユーザー情報の取得に失敗しました"})
			return
		}

		userInfo, err := utils.CurrentUser(c, redisConn, queries)
		if err != nil {
			log.Printf("Failed to retrieve current user: %v", err)
			c.Redirect(303, "/login")
			return
		}

		// ユーザーのブックマークしたツイートを取得
		tweets, err := queries.GetBookmarks(c, userID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "内部エラーが発生しました"})
			return
		}

		c.JSON(200, gin.H{
			"userId": userInfo.ID,
			"tweets": tweets,
		})
	}
}
