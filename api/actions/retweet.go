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

type ReTweetForm struct {
	UserId int64 `json:"userId"`
}

func RetweetAction(cfg config.Config, queries *sqlc.Queries) gin.HandlerFunc {
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

		// リツイートの存在確認
		retweet, err := queries.GetRetweet(c, sqlc.GetRetweetParams{
			OriginalTweetID: tweetID,
			UserID:          userID,
		})

		if err == sql.ErrNoRows {
			// リツイートが存在しない場合は追加
			// 元のツイートを取得
			originalTweet, err := queries.GetTweetByID(c, tweetID)
			if err != nil {
				c.JSON(500, gin.H{"error": "元のツイートの取得に失敗しました"})
				return
			}
			// リツイート内容をデータベースに保存
			insertedTweet, err := queries.InsertTweet(c, sqlc.InsertTweetParams{
				UserID:    originalTweet.UserID,
				Content:   originalTweet.Content,
				ImagePath: originalTweet.ImagePath,
				IsRetweet: true,
			})
			if err != nil {
				c.JSON(500, gin.H{"error": "リツイートのツイート保存に失敗しました"})
				return
			}

			retweetID := insertedTweet.ID
			err = queries.CreateRetweet(c, sqlc.CreateRetweetParams{
				TweetID:         retweetID,
				UserID:          userID,
				OriginalTweetID: tweetID,
			})
			if err != nil {
				c.JSON(500, gin.H{"error": "リツイートの保存に失敗しました"})
				return
			}

			c.JSON(200, gin.H{"message": "リツイート完了"})
		} else if err == nil {
			// トランザクション開始
			tx, err := db.DbConn().Begin()
			if err != nil {
				c.JSON(500, gin.H{"error": "トランザクションの開始に失敗しました"})
				return
			}

			// リツイートが存在する場合、削除
			err = queries.DeleteRetweet(c, sqlc.DeleteRetweetParams{
				OriginalTweetID: retweet.OriginalTweetID,
				UserID:          retweet.UserID,
			})
			if err != nil {
				tx.Rollback()
				c.JSON(500, gin.H{"error": "リツイートの削除に失敗しました"})
				return
			}

			// 関連するツイートを削除
			err = queries.DeleteRetweetOfTweet(c, retweet.TweetID)
			if err != nil {
				tx.Rollback()
				c.JSON(500, gin.H{"error": "ツイートの削除に失敗しました"})
				return
			}

			// トランザクションのコミット
			if err := tx.Commit(); err != nil {
				c.JSON(500, gin.H{"error": "トランザクションのコミットに失敗しました"})
				return
			}
			c.JSON(200, gin.H{"message": "リツイートを取り消しました"})
		} else {
			c.JSON(500, gin.H{"error": "リツイートの取得中にエラーが発生しました"})
			return
		}
	}
}
