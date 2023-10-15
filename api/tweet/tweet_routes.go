package tweet

import (
	"example.com/golang_twitter/api/middleware"
	"example.com/golang_twitter/config"
	"github.com/gin-gonic/gin"
)

func TweetRoutes(router *gin.RouterGroup, cfg config.Config) {
	tweet := router.Group("/tweet")
	{
		tweet.POST("/post", middleware.AuthRequired(), postTweet(cfg))        // ツイート投稿機能
		tweet.GET("/tweets", middleware.AuthRequired(), GetTweetsAsJSON(cfg)) // ページネーション用のツイートリスト取得
	}
}
