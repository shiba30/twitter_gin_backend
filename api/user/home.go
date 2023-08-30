package user

import (
	"context"
	"log"

	"example.com/golang_twitter/api/interfaces"
	"example.com/golang_twitter/util"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
)

const DefaultAvatarImage = "/static/img/default_avatar.png"

type RedisGetter interface {
	Get(ctx context.Context, key string) *redis.StringCmd
}

// ホーム表示
func ShowHome(c *gin.Context, redisConn *interfaces.RedisConn) {
	userInfo, err := util.CurrentUser(c, redisConn)
	if err != nil {
		log.Printf("Failed to retrieve current user: %v", err)
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
