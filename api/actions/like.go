package actions

import (
	"database/sql"
	"log"
	"net/http"
	"strconv"

	"example.com/golang_twitter/config"
	"example.com/golang_twitter/db"
	sqlc "example.com/golang_twitter/db/sqlc"
	"example.com/golang_twitter/utils"
	"github.com/gin-gonic/gin"
)

func LikeAction(cfg config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		// ツイートIDの取得
		tweetID, err := strconv.ParseInt(c.Param("id"), 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid tweet ID"})
			return
		}

		// セッションからユーザー情報を取得
		sessionID, err := c.Cookie("session_id")
		if err != nil {
			log.Printf("Failed to retrieve session ID: %v", err)
			c.JSON(500, gin.H{"error": "セッションIDの取得に失敗しました"})
			return
		}
		userID, err := utils.GetSessionUserId(c, sessionID)
		if err != nil {
			log.Printf("Failed to retrieve userId from session: %v", err)
			c.JSON(500, gin.H{"error": "セッションからユーザー情報の取得に失敗しました"})
			return
		}

		// いいねのレコード確認
		queries := sqlc.New(db.DbConn())
		_, err = queries.GetLike(c, sqlc.GetLikeParams{
			TweetID: tweetID,
			UserID:  userID,
		})

		if err == sql.ErrNoRows {
			// いいねが存在しない場合は追加
			err = queries.CreateLike(c, sqlc.CreateLikeParams{
				TweetID: tweetID,
				UserID:  userID,
			})
			if err != nil {
				c.JSON(500, gin.H{"error": "いいねの保存に失敗しました"})
				return
			}
			c.JSON(200, gin.H{"message": "いいね完了"})
		} else if err == nil {
			// いいねが存在する場合、削除
			err = queries.DeleteLike(c, sqlc.DeleteLikeParams{
				TweetID: tweetID,
				UserID:  userID,
			})
			if err != nil {
				c.JSON(500, gin.H{"error": "いいねの削除に失敗しました"})
				return
			}
			c.JSON(200, gin.H{"message": "いいねを取り消しました"})
		} else {
			c.JSON(500, gin.H{"error": "いいねの確認中にエラーが発生しました"})
			return
		}
	}
}
