package tweet

import (
	"log"
	"strconv"

	"example.com/golang_twitter/api/interfaces"
	"example.com/golang_twitter/config"
	sqlc "example.com/golang_twitter/db/sqlc"
	"example.com/golang_twitter/utils"
	"github.com/gin-gonic/gin"
)

func GetTweetList(c *gin.Context, cfg config.Config, queries *sqlc.Queries, userId int64) ([]sqlc.GetTweetsRow, error) {

	// ページネーション設定
	page := c.DefaultQuery("page", "1")
	pageSize := c.DefaultQuery("pageSize", cfg.DefaultPageSize)

	pageNum, err := strconv.Atoi(page)
	if err != nil {
		log.Printf("failed to convert to integer: %v", err)
		c.JSON(500, gin.H{"error": "システムエラー"})
		return nil, err
	}

	pageSizeNum, err := strconv.Atoi(pageSize)
	if err != nil {
		log.Printf("failed to convert to integer: %v", err)
		c.JSON(500, gin.H{"error": "システムエラー"})
		return nil, err
	}

	// ツイート取得処理
	tweets, err := queries.GetTweets(c, sqlc.GetTweetsParams{
		Limit:  int32(pageSizeNum),
		Offset: int32((pageNum - 1) * pageSizeNum),
		UserID: userId,
	})
	if err != nil {
		log.Printf("Error retrieving tweets from the database: %v", err)
		c.JSON(500, gin.H{"error": "ツイートの取得に失敗しました"})
		return nil, err
	}
	log.Println("successful retrieving tweets from the database")

	return tweets, nil
}

func GetTweetsAsJSON(cfg config.Config, redisConn *interfaces.RedisConn, queries *sqlc.Queries) gin.HandlerFunc {
	return func(c *gin.Context) {
		userInfo, err := utils.CurrentUser(c, redisConn, queries)
		if err != nil {
			log.Printf("Failed to get userInfo: %v", err)
			c.JSON(500, gin.H{"error": "プロフィール情報の取得に失敗しました"})
			return
		}

		tweets, err := GetTweetList(c, cfg, queries, userInfo.ID)
		if err != nil {
			log.Printf("Failed to retrieve tweets: %v", err)
			c.JSON(500, gin.H{"error": "ツイートの取得に失敗しました"})
			return
		}

		c.JSON(200, gin.H{
			"tweets":   tweets,
			"userInfo": userInfo,
		})
	}
}
