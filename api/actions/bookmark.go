package actions

import (
	"database/sql"
	"log"
	"net/http"
	"strconv"

	"example.com/golang_twitter/config"
	sqlc "example.com/golang_twitter/db/sqlc"
	"example.com/golang_twitter/utils"
	"github.com/gin-gonic/gin"
)

func PostBookmarkAction(cfg config.Config, queries *sqlc.Queries) gin.HandlerFunc {
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

		// ブックマークのレコード確認
		_, err = queries.GetBookmark(c, sqlc.GetBookmarkParams{
			TweetID: tweetID,
			UserID:  userID,
		})

		if err == sql.ErrNoRows {
			// ブックマークが存在しない場合は追加
			err = queries.CreateBookmark(c, sqlc.CreateBookmarkParams{
				TweetID: tweetID,
				UserID:  userID,
			})
			if err != nil {
				c.JSON(500, gin.H{"error": "ブックマークの保存に失敗しました"})
				return
			}
			c.JSON(200, gin.H{"message": "ブックマーク完了"})
		} else if err == nil {
			// ブックマークが存在する場合、削除
			err = queries.DeleteBookmark(c, sqlc.DeleteBookmarkParams{
				TweetID: tweetID,
				UserID:  userID,
			})
			if err != nil {
				c.JSON(500, gin.H{"error": "ブックマークの削除に失敗しました"})
				return
			}
			c.JSON(200, gin.H{"message": "ブックマークを取り消しました"})
		} else {
			c.JSON(500, gin.H{"error": "ブックマークの確認中にエラーが発生しました"})
			return
		}
	}
}
