package util

import (
	"fmt"

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
