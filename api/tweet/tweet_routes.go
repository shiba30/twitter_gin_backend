package tweet

import (
	"example.com/golang_twitter/api/middleware"
	"example.com/golang_twitter/api/reply"
	"example.com/golang_twitter/config"
	"github.com/gin-gonic/gin"
)

func TweetRoutes(router *gin.RouterGroup, cfg config.Config) {
	tweet := router.Group("/tweets")
	{
		tweet.POST("", middleware.AuthRequired(), postTweet(cfg))         // ツイート投稿機能
		tweet.GET("", middleware.AuthRequired(), GetTweetsAsJSON(cfg))    // ページネーション用のツイートリスト取得機能
		tweet.GET("/:id", middleware.AuthRequired(), GetTweetDetail(cfg)) // ツイート詳細取得機能

		// リプライ関連のルーティング
		tweet.POST("/:id/replies", middleware.AuthRequired(), reply.PostReply(cfg)) // コメント投稿機能
		tweet.GET("/:id/replies", middleware.AuthRequired(), reply.GetReplies(cfg)) // コメント一覧取得機能
	}
}
