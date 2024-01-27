package tweet

import (
	"log"
	"strconv"

	"example.com/golang_twitter/config"
	"example.com/golang_twitter/db"
	sqlc "example.com/golang_twitter/db/sqlc"
	"example.com/golang_twitter/utils"
	"github.com/gin-gonic/gin"
)

func GetTweetList(c *gin.Context, cfg config.Config, userId int64) ([]sqlc.GetTweetsRow, error) {

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
		Limit:      int32(pageSizeNum),
		Offset:     int32((pageNum - 1) * pageSizeNum),
		FollowerID: userId,
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
		// セッションからユーザー情報を取得
		sessionID, err := c.Cookie("session_id")
		if err != nil {
			log.Printf("Failed to retrieve session ID: %v", err)
			c.JSON(500, gin.H{"error": "セッションIDの取得に失敗しました"})
			return
		}
		currentUserId, err := utils.GetSessionUserId(c, sessionID)
		if err != nil {
			log.Printf("Failed to retrieve userId from session: %v", err)
			c.JSON(500, gin.H{"error": "セッションからユーザー情報の取得に失敗しました"})
			return
		}

		tweets, err := GetTweetList(c, cfg, currentUserId)
		if err != nil {
			log.Printf("Failed to retrieve tweets: %v", err)
			c.JSON(500, gin.H{"error": "ツイートの取得に失敗しました"})
			return
		}

		c.JSON(200, gin.H{
			"tweets":        tweets,
			"currentUserId": currentUserId,
		})
	}
}
