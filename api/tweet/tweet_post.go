package tweet

import (
	"database/sql"
	"encoding/base64"
	"fmt"
	"log"
	"os"
	"time"

	"example.com/golang_twitter/api/middleware"
	"example.com/golang_twitter/db"
	sqlc "example.com/golang_twitter/db/sqlc"
	"example.com/golang_twitter/util"
	"github.com/gin-gonic/gin"
)

const dir = "uploaded_images/"

type TweetForm struct {
	ID      int64  `json:"id"`
	UserId  int64  `json:"userId"`
	Content string `json:"content"`
	Image   string `json:"image,omitempty"`
}

func TweetRoutes(router *gin.RouterGroup) {
	tweet := router.Group("/tweet")
	{
		tweet.POST("/post", middleware.AuthRequired(), postTweet)
	}
}

func postTweet(c *gin.Context) {
	var form TweetForm

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
	userId, err := util.GetSessionUserId(c, sessionID)
	if err != nil {
		log.Printf("Failed to retrieve userId from session: %v", err)
		c.JSON(500, gin.H{"error": "セッションからユーザー情報の取得に失敗しました"})
		return
	}

	// Base64画像データのデコード (もし画像がある場合)
	var imagePath sql.NullString
	if form.Image != "" {
		var decodeErr error
		imageData, decodeErr := base64.StdEncoding.DecodeString(form.Image)
		if decodeErr != nil {
			log.Printf("failed to decode base64 image data: %v", decodeErr)
			c.JSON(400, gin.H{"error": "画像のデコードに失敗しました"})
			return
		}
		if _, err := os.Stat(dir); os.IsNotExist(err) {
			// ディレクトリが存在しない場合、作成します
			os.MkdirAll(dir, 0755)
		}
		// 現在のタイムスタンプを使用して、一意のファイル名を生成
		filename := fmt.Sprintf("%d_%d.png", time.Now().Unix(), form.UserId)
		path := dir + filename
		imagePath = sql.NullString{String: path, Valid: true}

		// 画像データをローカルのディレクトリに保存
		err := os.WriteFile(path, imageData, 0644)
		if err != nil {
			log.Printf("failed to save image data to local file: %v", err)
			c.JSON(500, gin.H{"error": "画像の保存に失敗しました"})
			return
		}
	} else {
		imagePath = sql.NullString{Valid: false}
	}

	// ツイート保存処理
	// sqlcのQueriesオブジェクトを初期化
	queries := sqlc.New(db.DbConn())
	_, err = queries.InsertTweet(c, sqlc.InsertTweetParams{
		UserID:    userId,
		Content:   form.Content,
		ImagePath: imagePath,
	})
	if err != nil {
		log.Printf("failed to save tweet: %v", err)
		c.JSON(500, gin.H{"error": "ツイート保存に失敗しました"})
		return
	}

	log.Println("successful tweet")

	c.JSON(200, gin.H{"user_id": form.UserId, "content": form.Content, "image_path": imagePath})

}
