package reply

import (
	"database/sql"
	"encoding/base64"
	"fmt"
	"log"
	"os"
	"time"

	"example.com/golang_twitter/config"
	"example.com/golang_twitter/db"
	sqlc "example.com/golang_twitter/db/sqlc"
	"example.com/golang_twitter/util"
	"github.com/gin-gonic/gin"
)

type ReplyForm struct {
	ReplyTo int64  `json:"replyTo"`
	UserId  int64  `json:"userId"`
	Content string `json:"content"`
	Image   string `json:"image,omitempty"`
}

func postReply(cfg config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		form := ReplyForm{}
		dir := cfg.UploadedImagesDir

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

		log.Println(newNullInt64(form.ReplyTo))
		log.Println(userId)
		log.Println(form.Content)
		log.Println(imagePath)

		// コメントデータの保存処理
		queries := sqlc.New(db.DbConn())
		_, err = queries.InsertReply(c, sqlc.InsertReplyParams{
			ReplyTo:   newNullInt64(form.ReplyTo),
			UserID:    userId,
			Content:   form.Content,
			ImagePath: imagePath,
		})
		if err != nil {
			log.Printf("failed to save reply: %v", err)
			c.JSON(500, gin.H{"error": "リプライ保存に失敗しました"})
			return
		}

		log.Println("successful reply")

		// 成功レスポンス
		c.JSON(200, gin.H{"message": "Comment posted successfully"})
	}
}

func newNullInt64(i int64) sql.NullInt64 {
	if i == 0 {
		return sql.NullInt64{Valid: false}
	}
	return sql.NullInt64{Int64: i, Valid: true}
}
