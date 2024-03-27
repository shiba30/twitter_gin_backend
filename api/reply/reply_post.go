package reply

import (
	"log"
	"strconv"

	"example.com/golang_twitter/api/interfaces"
	"example.com/golang_twitter/config"
	sqlc "example.com/golang_twitter/db/sqlc"
	"example.com/golang_twitter/utils"
	"github.com/gin-gonic/gin"
)

type ReplyForm struct {
	Content string `json:"content"`
	Image   string `json:"image,omitempty"`
}

func PostReply(cfg config.Config, redisConn *interfaces.RedisConn, queries *sqlc.Queries) gin.HandlerFunc {
	return func(c *gin.Context) {
		form := ReplyForm{}

		// リクエストデータの確認
		if err := c.ShouldBind(&form); err != nil {
			log.Printf("failed to bind tweet data: %v", err)
			c.JSON(400, gin.H{"error": "リプライに失敗しました"})
			return
		}

		// 140文字以内の制限
		if len(form.Content) > 140 {
			log.Printf("exceeds 140 characters: %v", form.Content)
			c.JSON(400, gin.H{"error": "ツイートは140文字以内にしてください"})
			return
		}

		tweetID, err := strconv.ParseInt(c.Param("id"), 10, 64)
		if err != nil {
			c.JSON(400, gin.H{"error": "ツイートIDが不正です"})
			return
		}

		// ユーザ情報の取得
		userInfo, err := utils.CurrentUser(c, redisConn, queries)
		if err != nil {
			c.JSON(500, gin.H{"error": "プロフィール情報の取得に失敗しました"})
			return
		}

		// 画像処理部分を共通関数に置き換え
		imagePath, err := utils.ProcessImage(form.Image, cfg.UploadedImagesDir, userInfo.ID)
		if err != nil {
			c.JSON(500, gin.H{"error": "画像の処理に失敗しました"})
			return
		}

		// コメントデータの保存処理
		_, err = queries.InsertReply(c, sqlc.InsertReplyParams{
			TweetID:   tweetID,
			UserID:    userInfo.ID,
			Content:   form.Content,
			ImagePath: imagePath,
		})
		if err != nil {
			log.Printf("failed to save reply: %v", err)
			c.JSON(500, gin.H{"error": "リプライ保存に失敗しました"})
			return
		}

		// 通知テーブルにコメントを登録
		tweet, err := queries.GetTweetByID(c, tweetID)
		if err != nil {
			c.JSON(500, gin.H{"error": "ツイート情報の取得に失敗しました"})
			return
		}
		if tweet.UserID != userInfo.ID {
			_, err = queries.InsertNotification(c, sqlc.InsertNotificationParams{
				UserID:      tweet.UserID,
				ActionType:  "comment",
				ReferenceID: tweetID,
			})
			if err != nil {
				c.JSON(500, gin.H{"error": "通知の登録に失敗しました"})
				return
			}
		}

		log.Println("successful reply")

		// 成功レスポンス
		c.JSON(200, gin.H{"message": "Comment posted successfully"})
	}
}
