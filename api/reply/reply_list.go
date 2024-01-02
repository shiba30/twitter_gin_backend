package reply

import (
	"database/sql"
	"log"
	"net/http"
	"strconv"

	"example.com/golang_twitter/config"
	"example.com/golang_twitter/db"
	sqlc "example.com/golang_twitter/db/sqlc"
	"github.com/gin-gonic/gin"
)

func getReply(cfg config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		// ツイートIDの取得
		tweetID, err := strconv.ParseInt(c.Param("id"), 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid tweet ID"})
			return
		}

		// リプライデータの取得処理
		queries := sqlc.New(db.DbConn())
		replies, err := queries.GetTweetDetailReply(c, sql.NullInt64{Int64: tweetID, Valid: true})
		if err != nil {
			log.Printf("failed to get replies: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "リプライの取得に失敗しました"})
			return
		}

		// リプライデータのレスポンス
		c.JSON(http.StatusOK, gin.H{
			"replies": replies,
		})
	}
}
