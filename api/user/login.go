package user

import (
	"log"

	"example.com/golang_twitter/api/interfaces"
	"example.com/golang_twitter/config"
	db "example.com/golang_twitter/db"
	sqlc "example.com/golang_twitter/db/sqlc"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type loginForm struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func LoginRoutes(router *gin.RouterGroup, cfg config.Config, redisConn *interfaces.RedisConn) {
	user := router.Group("/user")
	{
		user.POST("/login", func(c *gin.Context) {
			login(c, cfg, redisConn)
		})
	}
}

// ログイン機能
func login(c *gin.Context, cfg config.Config, redisConn *interfaces.RedisConn) {
	form := loginForm{}

	// リクエストデータの確認
	if err := c.ShouldBindJSON(&form); err != nil {
		log.Printf("failed to bind form data: %v", err)
		c.JSON(401, gin.H{"error": "ログイン認証に失敗しました"})
		return
	}

	// ユーザ情報取得
	// sqlcのQueriesオブジェクトを初期化
	queries := sqlc.New(db.DbConn())
	userInfo, err := queries.GetUserByEmail(c, form.Email)
	if err != nil {
		log.Printf("failed to get user info: %v", err)
		c.JSON(401, gin.H{"error": "ログイン認証に失敗しました"})
		return
	}

	// ユーザアクティブ確認
	if !userInfo.IsActive {
		log.Printf("user not active : %v", err)
		c.JSON(401, gin.H{"error": "ログイン認証に失敗しました"})
		return
	}

	// パスワード比較
	err = bcrypt.CompareHashAndPassword([]byte(userInfo.Password), []byte(form.Password))
	if err != nil {
		log.Printf("password do not match: %v", err)
		c.JSON(401, gin.H{"error": "ログイン認証に失敗しました"})
		return
	}

	// UUIDを生成して、セッションIDとして使用
	sessionID := uuid.New().String()

	// redisにsession情報を保存
	// セッション有効期限:0(無期限)
	err = redisConn.SetSession(c, sessionID, userInfo.ID, 0)
	if err != nil {
		log.Printf("failed to set session information: %v", err)
		c.JSON(500, gin.H{"error": "内部エラーが発生しました"})
		return
	}

	// Set-Cookie headerを追加し、レスポンス返却
	c.SetCookie("session_id", sessionID, 3600, "/", "", false, true)
	c.JSON(200, gin.H{"status": "ログインに成功しました"})

}
