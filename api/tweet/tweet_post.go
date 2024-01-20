package tweet

import (
	"log"

	"example.com/golang_twitter/config"
	"example.com/golang_twitter/db"
	sqlc "example.com/golang_twitter/db/sqlc"
	"example.com/golang_twitter/utils"
	"github.com/gin-gonic/gin"
)

type TweetForm struct {
	ID      int64  `json:"id"`
	UserId  int64  `json:"userId"`
	Content string `json:"content"`
	Image   string `json:"image,omitempty"`
}

func postTweet(cfg config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		form := TweetForm{}

		// リクエストデータの確認
		if err := c.ShouldBindJSON(&form); err != nil {
			log.Printf("failed to bind tweet data: %v", err)
			c.JSON(400, gin.H{"error": "ツイートに失敗しました"})
			return
		}

		// 140文字以内の制限
		if len(form.Content) > 140 {
			log.Printf("exceeds 140 characters: %v", form.Content)
			c.JSON(400, gin.H{"error": "ツイートは140文字以内にしてください"})
			return
		}

		// セッションからユーザー情報を取得
		sessionID, err := c.Cookie("session_id")
		if err != nil {
			log.Printf("Failed to retrieve session ID: %v", err)
			c.JSON(500, gin.H{"error": "セッションIDの取得に失敗しました"})
			return
		}
		userId, err := utils.GetSessionUserId(c, sessionID)
		if err != nil {
			log.Printf("Failed to retrieve userId from session: %v", err)
			c.JSON(500, gin.H{"error": "セッションからユーザー情報の取得に失敗しました"})
			return
		}

		// 画像処理部分を共通関数に置き換え
		imagePath, err := utils.ProcessImage(form.Image, cfg.UploadedImagesDir, userId)
		if err != nil {
			c.JSON(500, gin.H{"error": "画像の処理に失敗しました"})
			return
		}

		// ツイート保存処理
		// sqlcのQueriesオブジェクトを初期化
		queries := sqlc.New(db.DbConn())
		_, err = queries.InsertTweet(c, sqlc.InsertTweetParams{
			UserID:    userId,
			Content:   form.Content,
			ImagePath: imagePath,
			IsRetweet: false,
		})
		if err != nil {
			log.Printf("failed to save tweet: %v", err)
			c.JSON(500, gin.H{"error": "ツイート保存に失敗しました"})
			return
		}

		log.Println("successful tweet")

		c.JSON(200, gin.H{"user_id": form.UserId, "content": form.Content, "image_path": imagePath})

	}
}
