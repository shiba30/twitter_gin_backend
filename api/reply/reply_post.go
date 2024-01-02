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

type ReplyForm struct {
	UserId  int64          `form:"user_id"`
	Content string         `form:"content"`
	Image   sql.NullString `form:"image,omitempty"`
	ReplyTo int64          `form:"reply_to"`
}

func postReply(cfg config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		form := ReplyForm{}

		// リクエストデータの確認
		if err := c.ShouldBind(&form); err != nil {
			log.Printf("failed to bind tweet data: %v", err)
			c.JSON(400, gin.H{"error": "リプライに失敗しました"})
			return
		}

		// ツイートIDの取得
		tweetId, err := strconv.ParseInt(c.Param("id"), 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid tweet ID"})
			return
		}

		// コメントデータの保存処理
		queries := sqlc.New(db.DbConn())
		_, err = queries.InsertReply(c, sqlc.InsertReplyParams{
			UserID:    form.UserId,
			Content:   form.Content,
			ImagePath: form.Image,
			ReplyTo:   newNullInt64(tweetId),
		})
		if err != nil {
			log.Printf("failed to save reply: %v", err)
			c.JSON(500, gin.H{"error": "リプライ保存に失敗しました"})
			return
		}

		log.Println("successful commnet")

		// 成功レスポンス
		c.JSON(200, gin.H{"message": "Comment posted successfully"})
	}
}

func newNullInt64(i int64) sql.NullInt64 {
	return sql.NullInt64{Int64: i, Valid: i != 0}
}
