package reply

import (
	"log"
	"net/http"
	"strconv"

	"example.com/golang_twitter/api/interfaces"
	"example.com/golang_twitter/config"
	sqlc "example.com/golang_twitter/db/sqlc"
	"github.com/gin-gonic/gin"
)

func DeleteReply(cfg config.Config, redisConn *interfaces.RedisConn, queries *sqlc.Queries) gin.HandlerFunc {
	return func(c *gin.Context) {
		// リプライIDの取得
		replyID, err := strconv.ParseInt(c.Param("id"), 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid reply ID"})
			return
		}

		// リプライデータの削除処理
		err = queries.DeleteReply(c, replyID)
		if err != nil {
			log.Printf("failed to delete reply: %v", err)
			c.JSON(500, gin.H{"error": "リプライの削除失敗しました"})
			return
		}

		log.Println("successful delete reply")

		// 成功レスポンス
		c.JSON(200, gin.H{"message": "Comment deleted successfully"})
	}
}
