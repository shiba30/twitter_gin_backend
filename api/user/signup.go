package user

import (
	"database/sql"
	"log"
	"net/mail"
	"regexp"

	"example.com/golang_twitter/config"
	db "example.com/golang_twitter/db"
	sqlc "example.com/golang_twitter/db/sqlc"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

type signupForm struct {
	Email       string `json:"email"`
	Password    string `json:"password"`
	DisplayName string `json:"displayName"`
}

func signupHandler(cfg config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		signup(c, cfg)
	}
}

func SignupRoutes(router *gin.RouterGroup, cfg config.Config) {
	user := router.Group("/user")
	{
		user.POST("/signup", signupHandler(cfg))
		user.GET("/verify/:token", activateUser) // 確認メールのリンクが踏まれた時に呼び出される関数
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
func signup(c *gin.Context, cfg config.Config) {
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
	_, err := queries.GetUserByEmail(c, form.Email)
	if err != nil {
		if err != sql.ErrNoRows {
			log.Printf("failed to check if user exists: %v", err)
			c.JSON(500, gin.H{"error": "サインアップに失敗しました"})
			return
		}
	} else {
		// 既にユーザが登録されている
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
	user, err := queries.CreateUser(c, sqlc.CreateUserParams{
		Email:       form.Email,
		Password:    string(hash_pwd),
		DisplayName: form.DisplayName,
	})
	if err != nil {
		log.Printf("failed to create user: %v", err)
		c.JSON(500, gin.H{"error": "ユーザ情報登録に失敗しました"})
		return
	}

	// アクティベーショントークン生成
	activationToken, err := generateActivationToken(cfg, user)
	if err != nil {
		log.Printf("failed to generate activation token: %v", err)
		c.JSON(500, gin.H{"error": "アクティベーショントークンの生成に失敗しました"})
		return
	}

	// アクティベーショントークンをDBに保存
	_, err = queries.UpdateUser(c, sqlc.UpdateUserParams{
		ID:              user.ID,
		ActivationToken: sql.NullString{String: activationToken, Valid: true},
		IsActive:        false,
	})
	if err != nil {
		log.Printf("failed to save activation token: %v", err)
		c.JSON(500, gin.H{"error": "アクティベーショントークンの保存に失敗しました"})
		return
	}

	// アクティベーションメール送信
	err = sendActivationEmail(cfg, user, activationToken)
	if err != nil {
		log.Printf("failed to send activation email: %v", err)
		c.JSON(500, gin.H{"error": "アクティベーションメールの送信に失敗しました"})
		return
	}

	log.Printf("created user: %v", user)
	c.JSON(200, gin.H{"status": "ユーザ情報登録に成功しました"})
}
