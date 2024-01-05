package utils

import (
	"database/sql"
	"encoding/base64"
	"fmt"
	"log"
	"os"
	"time"

	"example.com/golang_twitter/api/interfaces"
	"example.com/golang_twitter/db"
	sqlc "example.com/golang_twitter/db/sqlc"
	"github.com/gin-gonic/gin"
)

func CurrentUser(c *gin.Context, redisConn *interfaces.RedisConn) (*sqlc.GetUserInfoRow, error) {
	sessionID, _ := c.Cookie("session_id")

	// RedisからユーザIDを取得
	userId, err := redisConn.GetUserIdFromSession(c, sessionID)
	if err != nil {
		return nil, err
	}

	// データベースからユーザ情報を取得
	queries := sqlc.New(db.DbConn())
	userInfo, err := queries.GetUserInfo(c, userId)
	if err != nil {
		return nil, err
	}

	return &userInfo, nil
}

// セッションからuserIdを取得する関数
func GetSessionUserId(c *gin.Context, sessionID string) (int64, error) {
	// 保存されているRedisConnインスタンスを取得
	redisConn, ok := c.Get("redisConn")
	if !ok {
		return 0, fmt.Errorf("could not retrieve RedisConn from context")
	}

	rc, ok := redisConn.(*interfaces.RedisConn)
	if !ok {
		return 0, fmt.Errorf("unexpected type for RedisConn in context")
	}

	return rc.GetUserIdFromSession(c, sessionID)
}

// Base64エンコードされた画像データをデコードし、ディレクトリに保存して、パス返却
func ProcessImage(imageData string, dir string, userId int64) (sql.NullString, error) {
	var imagePath sql.NullString
	if imageData != "" {
		data, decodeErr := base64.StdEncoding.DecodeString(imageData)
		if decodeErr != nil {
			log.Printf("failed to decode base64 image data: %v", decodeErr)
			return imagePath, decodeErr
		}

		if _, err := os.Stat(dir); os.IsNotExist(err) {
			os.MkdirAll(dir, 0755)
		}

		filename := fmt.Sprintf("%d_%d.png", time.Now().Unix(), userId)
		path := dir + filename
		imagePath = sql.NullString{String: path, Valid: true}

		err := os.WriteFile(path, data, 0644)
		if err != nil {
			log.Printf("failed to save image data to local file: %v", err)
			return imagePath, err
		}
	} else {
		imagePath = sql.NullString{Valid: false}
	}
	return imagePath, nil
}
