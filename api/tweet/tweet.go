package tweet

import (
	"log"
	"net/http"
	"strconv"

	"example.com/golang_twitter/config"
	sqlc "example.com/golang_twitter/db/sqlc"
	"example.com/golang_twitter/utils"
	"github.com/gin-gonic/gin"
)

func GetTweetDetail(cfg config.Config, queries *sqlc.Queries) gin.HandlerFunc {
	return func(c *gin.Context) {
		tweetIDStr := c.Param("id") // URLからツイートIDを文字列として取得

		// 文字列のtweetIDをint64に変換
		tweetID, err := strconv.ParseInt(tweetIDStr, 10, 64)
		if err != nil {
			// tweetIDの変換に失敗した場合のエラー処理
			c.JSON(http.StatusBadRequest, gin.H{"error": "無効なツイートID"})
			return
		}

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

		// データベースからツイートの詳細情報を取得するロジック
		tweetDetail, err := queries.GetTweetDetail(c, sqlc.GetTweetDetailParams{
			ID:     tweetID,
			UserID: currentUserId,
		})
		if err != nil {
			// データベースからの取得に失敗した場合のエラー処理
			c.JSON(http.StatusInternalServerError, gin.H{"error": "ツイートの取得に失敗しました"})
			return
		}

		// ツイート詳細情報をクライアントに返す
		c.JSON(200, gin.H{
			"tweet_id":       tweetDetail.TweetID,
			"user_id":        tweetDetail.ID,
			"user_name":      tweetDetail.UserName,
			"profile_image":  tweetDetail.ProfileImage,
			"tweet_date":     tweetDetail.TweetDate,
			"tweet_content":  tweetDetail.TweetContent,
			"image_path":     tweetDetail.ImagePath,
			"replies_count":  tweetDetail.RepliesCount,
			"likes_count":    tweetDetail.LikesCount,
			"retweets_count": tweetDetail.RetweetsCount,
			"is_bookmarked":  tweetDetail.IsBookmarked,
		})
	}
}
