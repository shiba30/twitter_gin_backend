package tweet

import (
	"log"
	"net/http"
	"strconv"

	"example.com/golang_twitter/api/interfaces"
	"example.com/golang_twitter/config"
	sqlc "example.com/golang_twitter/db/sqlc"
	"example.com/golang_twitter/utils"
	"github.com/gin-gonic/gin"
)

func GetTweetDetail(cfg config.Config, redisConn *interfaces.RedisConn, queries *sqlc.Queries) gin.HandlerFunc {
	return func(c *gin.Context) {
		tweetIDStr := c.Param("id") // URLからツイートIDを文字列として取得

		// 文字列のtweetIDをint64に変換
		tweetID, err := strconv.ParseInt(tweetIDStr, 10, 64)
		if err != nil {
			// tweetIDの変換に失敗した場合のエラー処理
			c.JSON(http.StatusBadRequest, gin.H{"error": "無効なツイートID"})
			return
		}

		// ユーザ情報の取得
		userInfo, err := utils.CurrentUser(c, redisConn, queries)
		if err != nil {
			c.JSON(500, gin.H{"error": "プロフィール情報の取得に失敗しました"})
			return
		}

		// データベースからツイートの詳細情報を取得するロジック
		tweetDetail, err := queries.GetTweetDetail(c, sqlc.GetTweetDetailParams{
			ID:     tweetID,
			UserID: userInfo.ID,
		})
		if err != nil {
			// データベースからの取得に失敗した場合のエラー処理
			c.JSON(http.StatusInternalServerError, gin.H{"error": "ツイートの取得に失敗しました"})
			return
		}

		// ツイートのコメント情報を取得
		replies, err := queries.GetTweetDetailReply(c, tweetID)
		if err != nil {
			log.Printf("failed to get replies: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "リプライの取得に失敗しました"})
			return
		}

		// ツイート詳細情報とコメント情報をクライアントに返す
		c.JSON(200, gin.H{
			"tweetDetail":   tweetDetail,
			"replies":       replies,
			"currentUserId": userInfo.ID,
		})
	}
}
