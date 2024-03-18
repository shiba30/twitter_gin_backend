package tweet

import (
	"log"
	"strconv"

	"example.com/golang_twitter/api/interfaces"
	sqlc "example.com/golang_twitter/db/sqlc"
	"example.com/golang_twitter/utils"
	"github.com/gin-gonic/gin"
)

func deleteTweet(redisConn *interfaces.RedisConn, queries *sqlc.Queries) gin.HandlerFunc {
	return func(c *gin.Context) {
		tweetID, err := strconv.ParseInt(c.Param("id"), 10, 64)
		if err != nil {
			c.JSON(400, gin.H{"error": "ツイートIDが不正です"})
			return
		}

		// ユーザ情報の取得
		userInfo, err := utils.CurrentUser(c, redisConn, queries)
		if err != nil {
			c.JSON(500, gin.H{"error": "プロフィール情報の取得に失敗しました"})
			return
		}

		// ツイート削除処理
		err = queries.DeleteTweet(c, sqlc.DeleteTweetParams{
			ID:     tweetID,
			UserID: userInfo.ID,
		})
		if err != nil {
			log.Printf("failed to delete tweet: %v", err)
			c.JSON(500, gin.H{"error": "ツイート削除に失敗しました"})
			return
		}

		log.Println("successful tweet delete")

		c.JSON(200, gin.H{"message": "ツイートが削除されました"})

	}
}
