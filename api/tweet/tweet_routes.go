package tweet

import (
	"example.com/golang_twitter/api/actions"
	"example.com/golang_twitter/api/follow"
	"example.com/golang_twitter/api/middleware"
	"example.com/golang_twitter/api/reply"
	"example.com/golang_twitter/config"
	sqlc "example.com/golang_twitter/db/sqlc"
	"github.com/gin-gonic/gin"
)

func TweetRoutes(router *gin.RouterGroup, cfg config.Config, queries *sqlc.Queries) {
	tweet := router.Group("/tweets")
	{
		tweet.POST("", middleware.AuthRequired(), postTweet(cfg, queries))         // ツイート投稿機能
		tweet.GET("", middleware.AuthRequired(), GetTweetsAsJSON(cfg, queries))    // ページネーション用のツイートリスト取得機能
		tweet.GET("/:id", middleware.AuthRequired(), GetTweetDetail(cfg, queries)) // ツイート詳細取得機能

		tweet.POST("/:id/replies", middleware.AuthRequired(), reply.PostReply(cfg, queries)) // コメント投稿機能
		tweet.GET("/:id/replies", middleware.AuthRequired(), reply.GetReplies(cfg, queries)) // コメント一覧取得機能

		tweet.POST("/:id/like", middleware.AuthRequired(), actions.LikeAction(cfg, queries))             // いいね機能
		tweet.POST("/:id/retweet", middleware.AuthRequired(), actions.RetweetAction(cfg, queries))       // リツイート機能
		tweet.POST("/:id/bookmark", middleware.AuthRequired(), actions.PostBookmarkAction(cfg, queries)) // ブックマーク追加削除機能

		tweet.POST("/:id/follow", middleware.AuthRequired(), follow.Follow(cfg, queries))     // フォロー機能
		tweet.POST("/:id/unfollow", middleware.AuthRequired(), follow.UnFollow(cfg, queries)) // フォロー解除機能
	}
}
