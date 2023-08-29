package user

import (
	"context"
	"log"

	"example.com/golang_twitter/api/interfaces"
	db "example.com/golang_twitter/db"
	sqlc "example.com/golang_twitter/db/sqlc"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
)

const DefaultAvatarImage = "/static/img/default_avatar.png"

type RedisGetter interface {
	Get(ctx context.Context, key string) *redis.StringCmd
}

// ホーム表示
func ShowHome(c *gin.Context, redisConn *interfaces.RedisConn) {
	sessionID, err := c.Cookie("session_id")
	if err != nil {
		log.Printf("Not authenticated: %v", err)
		c.Redirect(303, "/api/user/login")
		return
	}

	// Redisからユーザのメールアドレスを取得
	var userEmail string
	userEmail, err = redisConn.GetSession(sessionID)
	if err != nil {
		log.Printf("Session not valid: %v", err)
		c.Redirect(303, "/api/user/login")
		return
	}

	// データベースからユーザ情報を取得
	queries := sqlc.New(db.DbConn())
	userInfo, err := queries.GetUserInfo(c, userEmail)
	if err != nil {
		log.Printf("Failed to retrieve user info: %v", err)
		c.Redirect(303, "/api/user/login")
		return
	}

	var avatarImage string
	if userInfo.AvatarImage.Valid {
		avatarImage = userInfo.AvatarImage.String
	} else {
		avatarImage = DefaultAvatarImage
	}

	c.HTML(200, "home.html", gin.H{
		"userId":      userInfo.ID,
		"displayName": userInfo.DisplayName,
		"avatarImage": avatarImage,
	})
}
