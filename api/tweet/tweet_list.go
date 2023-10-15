package tweet

import (
	"log"
	"strconv"

	"example.com/golang_twitter/config"
	"example.com/golang_twitter/db"
	sqlc "example.com/golang_twitter/db/sqlc"
	"github.com/gin-gonic/gin"
)

func GetTweetList(c *gin.Context, cfg config.Config) ([]sqlc.GetTweetsRow, error) {

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
	queries := sqlc.New(db.DbConn())
	tweets, err := queries.GetTweets(c, sqlc.GetTweetsParams{
		Limit:  int32(pageSizeNum),
		Offset: int32((pageNum - 1) * pageSizeNum),
	})
	if err != nil {
		log.Printf("Error retrieving tweets from the database: %v", err)
		c.JSON(500, gin.H{"error": "ツイートの取得に失敗しました"})
		return nil, err
	}
	log.Println("successful retrieving tweets from the database")

	return tweets, nil
}

func GetTweetsAsJSON(cfg config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		tweets, err := GetTweetList(c, cfg)
		if err != nil {
			log.Printf("Failed to retrieve tweets: %v", err)
			c.JSON(500, gin.H{"error": "ツイートの取得に失敗しました"})
			return
		}

		c.JSON(200, tweets)
	}
}
