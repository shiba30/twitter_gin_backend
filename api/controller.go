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
	sqlc "example.com/golang_twitter/db/sqlc"
	"github.com/gin-gonic/gin"
)

type Controller struct {
	Config    config.Config
	RedisConn *interfaces.RedisConn
	Queries   *sqlc.Queries
}

func NewController(cfg config.Config, redisConn *interfaces.RedisConn, queries *sqlc.Queries) *Controller {
	return &Controller{
		Config:    cfg,
		RedisConn: redisConn,
		Queries:   queries,
	}
}

func (c *Controller) Routes(router *gin.Engine) {
	// APIのルーティング
	api := router.Group("/")
	{
		user.SignupRoutes(api, c.Config, c.Queries)
		user.LoginRoutes(api, c.RedisConn, c.Queries)
		tweet.TweetRoutes(api, c.Config, c.Queries)
		bookmark.BookmarkRoutes(api, c.Config, c.RedisConn, c.Queries)
		message.MessageRoutes(api, c.Config, c.RedisConn, c.Queries)
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
		user.ShowHome(ctx, c.RedisConn, c.Config, c.Queries)
	})
}
