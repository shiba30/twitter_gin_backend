package tweet

import (
	"net/http"
	"strconv"

	"example.com/golang_twitter/config"
	"example.com/golang_twitter/constants"
	"example.com/golang_twitter/db"
	sqlc "example.com/golang_twitter/db/sqlc"
	"github.com/gin-gonic/gin"
)

func GetTweetDetail(cfg config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		tweetIDStr := c.Param("id") // URLからツイートIDを文字列として取得

		// 文字列のtweetIDをint64に変換
		tweetID, err := strconv.ParseInt(tweetIDStr, 10, 64)
		if err != nil {
			// tweetIDの変換に失敗した場合のエラー処理
			c.JSON(http.StatusBadRequest, gin.H{"error": "無効なツイートID"})
			return
		}

		// データベースからツイートの詳細情報を取得するロジック
		queries := sqlc.New(db.DbConn())
		tweetDetail, err := queries.GetTweetDetail(c, tweetID)
		if err != nil {
			// データベースからの取得に失敗した場合のエラー処理
			c.JSON(http.StatusInternalServerError, gin.H{"error": "ツイートの取得に失敗しました"})
			return
		}

		avatarImage := constants.DefaultAvatarImage
		if tweetDetail.AvatarImage.Valid {
			avatarImage = tweetDetail.AvatarImage.String
		}

		// ツイート詳細情報をクライアントに返す
		c.HTML(200, "tweet_detail.html", gin.H{
			"userId":       tweetDetail.ID,
			"displayName":  tweetDetail.DisplayName,
			"avatarImage":  avatarImage,
			"tweetId":      tweetDetail.TweetID,
			"UserImage":    tweetDetail.UserImage,
			"TweetContent": tweetDetail.TweetContent,
			"ImagePath":    tweetDetail.ImagePath,
			"tweetDate":    tweetDetail.TweetDate,
		})
	}
}
