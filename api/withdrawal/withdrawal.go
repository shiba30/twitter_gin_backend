package withdrawal

import (
	"log"

	"example.com/golang_twitter/api/interfaces"
	"example.com/golang_twitter/config"
	"example.com/golang_twitter/db"
	sqlc "example.com/golang_twitter/db/sqlc"
	"example.com/golang_twitter/utils"
	"github.com/gin-gonic/gin"
)

func WithdrawalRoutes(router *gin.RouterGroup, cfg config.Config, redisConn *interfaces.RedisConn, queries *sqlc.Queries) {
	router.DELETE("/withdrawal", Withdrawal(cfg, redisConn, queries))
}

func Withdrawal(cfg config.Config, redisConn *interfaces.RedisConn, queries *sqlc.Queries) gin.HandlerFunc {
	return func(c *gin.Context) {
		// ユーザID取得
		userInfo, err := utils.CurrentUser(c, redisConn, queries)
		if err != nil {
			log.Printf("Failed to get userInfo: %v", err)
			c.JSON(500, gin.H{"error": "プロフィール情報の取得に失敗しました"})
			return
		}

		tx, err := db.DbConn().Begin()
		if err != nil {
			c.JSON(500, gin.H{"error": "トランザクションの開始に失敗しました"})
			return
		}

		// ユーザーに関連するツイートの「コメント」を削除
		if err := queries.WithTx(tx).DeleteRepliesByUserID(c, userInfo.ID); err != nil {
			tx.Rollback()
			log.Printf("Failed to delete replies: %v", err)
			c.JSON(500, gin.H{"error": "返信の削除に失敗しました"})
			return
		}

		// ユーザーに関連するツイートの「いいね」を削除
		if err := queries.WithTx(tx).DeleteLikesByUserID(c, userInfo.ID); err != nil {
			tx.Rollback()
			log.Printf("Failed to delete likes: %v", err)
			c.JSON(500, gin.H{"error": "いいねの削除に失敗しました"})
			return
		}

		// ユーザーに関連するツイートの「リツイート」を削除
		if err := queries.WithTx(tx).DeleteRetweetsByUserID(c, userInfo.ID); err != nil {
			tx.Rollback()
			log.Printf("Failed to delete retweets: %v", err)
			c.JSON(500, gin.H{"error": "リツイートの削除に失敗しました"})
			return
		}

		// ユーザーに関連するツイートの「ブックマーク」を削除
		if err := queries.WithTx(tx).DeleteBookmarksByTweetUserID(c, userInfo.ID); err != nil {
			tx.Rollback()
			log.Printf("Failed to delete bookmarks: %v", err)
			c.JSON(500, gin.H{"error": "ブックマークの削除に失敗しました"})
			return
		}

		// ユーザーに関連するツイートの「ブックマーク」を削除
		if err := queries.WithTx(tx).DeleteBookmarksByUserID(c, userInfo.ID); err != nil {
			tx.Rollback()
			log.Printf("Failed to delete bookmarks: %v", err)
			c.JSON(500, gin.H{"error": "ブックマークの削除に失敗しました"})
			return
		}

		// ユーザーに関連するツイートを削除
		if err := queries.WithTx(tx).DeleteTweetsByUserID(c, userInfo.ID); err != nil {
			tx.Rollback()
			log.Printf("Failed to delete tweets: %v", err)
			c.JSON(500, gin.H{"error": "ツイートの削除に失敗しました"})
			return
		}

		// ユーザーに関連するツイートの「フォロー」を削除
		if err := queries.WithTx(tx).DeleteFollowsByUserID(c, userInfo.ID); err != nil {
			tx.Rollback()
			log.Printf("Failed to delete follows: %v", err)
			c.JSON(500, gin.H{"error": "フォローの削除に失敗しました"})
			return
		}

		// ユーザーに関連するツイートの「メッセージ」を削除
		if err := queries.WithTx(tx).DeleteMessagesByUserID(c, userInfo.ID); err != nil {
			tx.Rollback()
			log.Printf("Failed to delete messages: %v", err)
			c.JSON(500, gin.H{"error": "メッセージの削除に失敗しました"})
			return
		}

		// ユーザーに関連するツイートの「通知」を削除
		if err := queries.WithTx(tx).DeleteNotificationsByUserID(c, userInfo.ID); err != nil {
			tx.Rollback()
			log.Printf("Failed to delete notifications: %v", err)
			c.JSON(500, gin.H{"error": "通知の削除に失敗しました"})
			return
		}

		// 最後にユーザー情報を削除
		if err := queries.WithTx(tx).DeleteUser(c, userInfo.ID); err != nil {
			tx.Rollback()
			log.Printf("Failed to delete user: %v", err)
			c.JSON(500, gin.H{"error": "ユーザ削除処理に失敗しました"})
			return
		}

		tx.Commit()
		log.Printf("Deleted user: %v", userInfo.ID)
		c.JSON(200, gin.H{"message": "ユーザ削除が完了しました"})
	}
}
