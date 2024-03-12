package user

import (
	"context"
	"log"

	"example.com/golang_twitter/api/interfaces"
	"example.com/golang_twitter/api/tweet"
	"example.com/golang_twitter/config"
	sqlc "example.com/golang_twitter/db/sqlc"
	"example.com/golang_twitter/utils"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
)

type RedisGetter interface {
	Get(ctx context.Context, key string) *redis.StringCmd
}

// ホーム表示
func ShowHome(c *gin.Context, redisConn *interfaces.RedisConn, cfg config.Config, queries *sqlc.Queries) {
	userInfo, err := utils.CurrentUser(c, redisConn, queries)
	if err != nil {
		log.Printf("Failed to retrieve current user: %v", err)
		c.Redirect(303, "/login")
		return
	}

	tweets, err := tweet.GetTweetList(c, cfg, queries, userInfo.ID)
	if err != nil {
		log.Printf("Failed to retrieve tweets: %v", err)
		c.JSON(500, gin.H{"error": "ツイートの取得に失敗しました"})
		return
	}

	c.HTML(200, "home.html", gin.H{
		"userId":       userInfo.ID,
		"displayName":  userInfo.DisplayName,
		"profileImage": userInfo.ProfileImage,
		"tweets":       tweets,
	})
}
