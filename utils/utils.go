package utils

import (
	"crypto/rand"
	"database/sql"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"log"
	"os"
	"time"

	"example.com/golang_twitter/api/interfaces"
	sqlc "example.com/golang_twitter/db/sqlc"
	"github.com/gin-gonic/gin"
)

func CurrentUser(c *gin.Context, redisConn *interfaces.RedisConn, queries *sqlc.Queries) (*sqlc.GetUserInfoRow, error) {
	sessionID, _ := c.Cookie("session_id")

	// RedisからユーザIDを取得
	userId, err := redisConn.GetUserIdFromSession(c, sessionID)
	if err != nil {
		return nil, err
	}

	// データベースからユーザ情報を取得
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
			if err := os.MkdirAll(dir, 0755); err != nil {
				log.Printf("failed to create directory: %v", err)
				return imagePath, err
			}
		}

		randomPart, err := generateRandomString(8)
		if err != nil {
			log.Printf("failed to generate random part for filename: %v", err)
			return imagePath, err
		}

		filename := fmt.Sprintf("%d_%d_%s.png", time.Now().Unix(), userId, randomPart)
		path := dir + filename
		imagePath = sql.NullString{String: path, Valid: true}

		if err := os.WriteFile(path, data, 0644); err != nil {
			log.Printf("failed to save image data to local file: %v", err)
			return imagePath, err
		}
	} else {
		imagePath = sql.NullString{Valid: false}
	}
	return imagePath, nil
}

func generateRandomString(n int) (string, error) {
	b := make([]byte, n)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}
