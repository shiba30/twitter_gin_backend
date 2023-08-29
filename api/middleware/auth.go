package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

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
