/*
	全てのルーティングを行う
*/

package api

import (
	"example.com/golang_twitter/api/bookmark"
	"example.com/golang_twitter/api/interfaces"
	"example.com/golang_twitter/api/message"
	"example.com/golang_twitter/api/middleware"
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
	api := router.Group("/")
	{
		user.SignupRoutes(api, c.Config)
		user.LoginRoutes(api, c.Config, c.RedisConn)
		tweet.TweetRoutes(api, c.Config)
		bookmark.BookmarkRoutes(api, c.Config, c.RedisConn)
		message.MessageRoutes(api, c.Config, c.RedisConn)
	}

	// ページレンダリングのルーティング
	api.GET("/signup", func(c *gin.Context) {
		c.HTML(200, "signup.html", nil)
	})
	api.GET("/verification", func(c *gin.Context) {
		c.HTML(200, "verification.html", nil)
	})
	api.GET("/login", func(c *gin.Context) {
		c.HTML(200, "login.html", nil)
	})
	api.GET("/home", middleware.AuthRequired(), func(ctx *gin.Context) {
		user.ShowHome(ctx, c.RedisConn, c.Config)
	})
}
