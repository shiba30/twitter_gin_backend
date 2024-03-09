package user

import (
	"database/sql"
	"fmt"
	"log"
	"net/smtp"

	"example.com/golang_twitter/config"
	sqlc "example.com/golang_twitter/db/sqlc"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

// JWTで署名されたアクティベーショントークンを生成
func generateActivationToken(cfg config.Config, user sqlc.CreateUserRow) (string, error) {

	// トークンに含めるクレームを設定
	claims := jwt.MapClaims{}
	claims["authorized"] = true
	claims["user_id"] = user.ID

	// JWTを生成
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// 秘密鍵を使ってトークンに署名
	tokenString, err := token.SignedString([]byte(cfg.SecretKey))

	if err != nil {
		return "", err
	}

	return tokenString, nil
}

// アクティベーションメール送信
func sendActivationEmail(cfg config.Config, user sqlc.CreateUserRow, activationToken string) error {
	to := user.Email

	body := "To: " + to + "\r\n" +
		"Subject: メールアドレス確認依頼\r\n" +
		"\r\n" +
		"ユーザ登録が完了しました!\r\n" +
		"下記のリンクをクリックしてメールアドレスを確認してください。\r\n" +
		cfg.ClientAddress + "/verification/" + activationToken + "\r\n"

	err := smtp.SendMail(cfg.SmtpHost+":"+cfg.SmtpPort, nil, cfg.From, []string{to}, []byte(body))
	if err != nil {
		return fmt.Errorf("failed to send activation email: %v", err)
	}

	log.Print("confirmation email sent")
	return nil
}

// ユーザーをアクティブにする
func activateUser(c *gin.Context, queries *sqlc.Queries) {
	tokenStr := c.Param("token")
	token := sql.NullString{
		String: tokenStr,
		Valid:  tokenStr != "",
	}

	// アクティベーショントークンでユーザ検索
	userId, err := queries.GetUserByActivationToken(c, token)
	if err != nil {
		log.Printf("failed to find user by activation token: %v", err)
		c.JSON(400, gin.H{"error": "アクティベーショントークンが無効です"})
		return
	}

	// ユーザをアクティブ状態に更新
	_, err = queries.UpdateUser(c, sqlc.UpdateUserParams{
		ID:              userId,
		ActivationToken: sql.NullString{Valid: false},
		IsActive:        true,
	})
	if err != nil {
		log.Printf("failed to activate user: %v", err)
		c.JSON(500, gin.H{"error": "ユーザーのアクティベーションに失敗しました"})
	} else {
		log.Printf("activated user: %v", userId)
		c.JSON(200, gin.H{"status": "ユーザーのアクティベーションに成功しました"})
	}
}
