package user

import (
	"context"
	"log"
	"net/http"

	"example.com/golang_twitter/config"
	db "example.com/golang_twitter/db"
	sqlc "example.com/golang_twitter/db/sqlc"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type loginForm struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func LoginRoutes(router *gin.RouterGroup, cfg config.Config) {
	user := router.Group("/user")
	{
		user.GET("/login", func(c *gin.Context) {
			c.HTML(http.StatusOK, "login.html", nil)
		})
		user.POST("/login", func(c *gin.Context) {
			login(c)
		})
		user.GET("/home", AuthRequired(), func(c *gin.Context) {
			c.HTML(http.StatusOK, "home.html", nil)
		})
	}
}

func AuthRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		sessionID, err := c.Cookie("session_id")
		if err != nil || sessionID == "" {
			// ログイン認証していない場合、login.htmlにリダイレクト
			c.Redirect(http.StatusSeeOther, "/api/user/login")
			c.Abort()
			return
		}
		c.Next()
	}
}

var conn *redis.Client

func init() {
	conn = redis.NewClient(&redis.Options{
		Addr:     "redis:6379",
		Password: "",
		DB:       0,
	})
}

// ログイン機能
func login(c *gin.Context) {
	var form loginForm

	// リクエストデータの確認
	if err := c.ShouldBindJSON(&form); err != nil {
		log.Printf("failed to bind form data: %v", err)
		c.JSON(401, gin.H{"error": "ログイン認証に失敗しました"})
		return
	}

	// sqlcのQueriesオブジェクトを初期化
	queries := sqlc.New(db.DbConn())

	// ユーザ情報取得
	log.Printf("Email to find: %s", form.Email)

	userInfo, err := queries.GetUserInfo(context.Background(), form.Email)
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
	err = conn.Set(context.Background(), sessionID, userInfo.ID, 0).Err()
	if err != nil {
		log.Printf("failed to set session information: %v", err)
		c.JSON(500, gin.H{"error": "内部エラーが発生しました"})
		return
	}

	// Set-Cookie headerを追加し、レスポンス返却
	c.SetCookie("session_id", sessionID, 3600, "/", "", false, true)
	c.JSON(200, gin.H{"status": "ログインに成功しました"})

}
