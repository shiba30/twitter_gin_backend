package user

import (
	"context"
	"database/sql"
	"log"
	"net/http"
	"net/mail"
	"regexp"

	db "example.com/golang_twitter/db"
	sqlc "example.com/golang_twitter/db/sqlc"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

type signupForm struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func Routes(router *gin.RouterGroup) {
	user := router.Group("/user")
	{
		user.GET("/signup", func(c *gin.Context) {
			c.HTML(http.StatusOK, "signup.html", nil)
		})
		user.POST("/signup", signup)
		user.GET("/verification", func(c *gin.Context) {
			c.HTML(http.StatusOK, "verification.html", nil)
		})
	}
}

// パスワードバリデーション
func validatePassword(pwd string) bool {
	// 各正規表現を格納する配列
	passwordRegexps := []string{
		"^.{8,}$",     // 8文字指定
		"[a-zA-Z0-9]", // 半角英数字
		"[a-z]",       // 小文字が混合
		"[A-Z]",       // 大文字が混合
		"[!?-_]"}      // 対象の記号が1文字以上含まれる

	for _, reg := range passwordRegexps {
		if matched := regexp.MustCompile(reg).Match([]byte(pwd)); !matched {
			return false
		}
	}
	return true
}

// サインアップ機能
func signup(c *gin.Context) {
	var form signupForm

	// フォーム値チェック
	if err := c.BindJSON(&form); err != nil {
		log.Println("form value check error")
		c.JSON(400, gin.H{"error": "正しい情報を入力してください"})
		return
	}

	// メールアドレスチェック
	if _, err := mail.ParseAddress(form.Email); err != nil {
		log.Println("email check error")
		c.JSON(400, gin.H{"error": "正しいメールアドレスを入力してください"})
		return
	}

	// パスワードチェック
	if !(validatePassword(form.Password)) {
		log.Println("password check error")
		c.JSON(400, gin.H{"error": "入力されたパスワードは登録できません。\n[半角英数字]、[小・大文字]・[!?-_]の記号1文字以上含めた8文字にしてください"})
		return
	}

	// sqlcのQueriesオブジェクトを初期化
	queries := sqlc.New(db.DbConn())

	// 既に登録されているユーザであるかをチェック
	existingUser, err := queries.GetUserByEmail(context.Background(), form.Email)
	if err != nil && err != sql.ErrNoRows {
		log.Printf("failed to check if user exists: %v", err)
		c.JSON(500, gin.H{"error": "サインアップに失敗しました"})
		return
	}
	if existingUser != nil {
		c.JSON(400, gin.H{"error": "既に登録されています"})
		return
	}

	// パスワードをハッシュ化
	hash_pwd, err := bcrypt.GenerateFromPassword([]byte(form.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Printf("failed to hash password: %v", err)
		c.JSON(500, gin.H{"error": "ユーザ情報登録に失敗しました"})
		return
	}

	// ユーザーを作成
	user, err := queries.CreateUser(context.Background(), sqlc.CreateUserParams{
		Email:    form.Email,
		Password: string(hash_pwd),
	})
	if err != nil {
		log.Printf("failed to create user: %v", err)
		c.JSON(500, gin.H{"error": "ユーザ情報登録に失敗しました"})
		return
	}

	log.Printf("created user: %v", user)
	c.JSON(200, gin.H{"status": "ユーザ情報登録に成功しました"})
}
