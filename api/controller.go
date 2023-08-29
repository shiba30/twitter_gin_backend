/*
	全てのルーティングを行う
*/

package api

import (
	"net/http"

	"example.com/golang_twitter/api/interfaces"
	"example.com/golang_twitter/api/tweet"
	"example.com/golang_twitter/api/user"
	"example.com/golang_twitter/config"
	"github.com/gin-gonic/gin"
)

type Controller struct {
	Config    config.Config
	RedisConn *interfaces.RedisConn
}

func NewController(cfg config.Config, redisConn *interfaces.RedisConn) *Controller {
	return &Controller{
		Config:    cfg,
		RedisConn: redisConn,
	}
}

func (c *Controller) Routes(router *gin.Engine) {
	// APIのルーティング
	api := router.Group("/api")
	{
		user.SignupRoutes(api, c.Config)
		user.LoginRoutes(api, c.Config, c.RedisConn)
		tweet.TweetRoutes(api)
	}

	// ページレンダリングのルーティング
	api.GET("/user/signup", func(c *gin.Context) {
		c.HTML(http.StatusOK, "signup.html", nil)
	})
	api.GET("/user/verification", func(c *gin.Context) {
		c.HTML(http.StatusOK, "verification.html", nil)
	})
	api.GET("/user/login", func(c *gin.Context) {
		c.HTML(http.StatusOK, "login.html", nil)
	})
	api.GET("/user/home", func(ctx *gin.Context) {
		user.ShowHome(ctx, c.RedisConn)
	})
}
