package follow

import (
	"log"

	"example.com/golang_twitter/api/interfaces"
	sqlc "example.com/golang_twitter/db/sqlc"
	"example.com/golang_twitter/utils"
	"github.com/gin-gonic/gin"
)

type followForm struct {
	UserId int64 `json:"user_id"`
}

func followAction(c *gin.Context, isFollow bool, redisConn *interfaces.RedisConn, queries *sqlc.Queries) {
	form := followForm{}
	if err := c.ShouldBind(&form); err != nil {
		log.Printf("failed to bind data: %v", err)
		c.JSON(400, gin.H{"error": "処理に失敗しました"})
		return
	}

	// ログインユーザ情報の取得
	userInfo, err := utils.CurrentUser(c, redisConn, queries)
	if err != nil {
		c.JSON(500, gin.H{"error": "ログインユーザ情報の取得に失敗しました"})
		return
	}

	if isFollow {
		// フォロー
		err = queries.CreateFollow(c, sqlc.CreateFollowParams{
			FollowerID: userInfo.ID,
			FolloweeID: form.UserId,
		})
	} else {
		// フォロー解除
		err = queries.DeleteFollow(c, sqlc.DeleteFollowParams{
			FollowerID: userInfo.ID,
			FolloweeID: form.UserId,
		})
	}
	if err != nil {
		log.Printf("Error in follow action: %v", err)
		c.JSON(500, gin.H{"error": "処理に失敗しました"})
		return
	}

	// フォロー情報の取得
	followCount, err := queries.GetFollowByUserId(c, sqlc.GetFollowByUserIdParams{
		FollowerID: userInfo.ID,
		FolloweeID: form.UserId,
	})
	if err != nil {
		log.Printf("Failed to retrieve following: %v", err)
		c.JSON(500, gin.H{"error": "フォロー情報の取得に失敗しました"})
		return
	}

	log.Printf("successful follow execute")
	c.JSON(200, gin.H{"followCount": followCount})
}

// フォロー機能
func Follow(redisConn *interfaces.RedisConn, queries *sqlc.Queries) gin.HandlerFunc {
	return func(c *gin.Context) {
		followAction(c, true, redisConn, queries)
	}
}

// フォロー解除機能
func UnFollow(redisConn *interfaces.RedisConn, queries *sqlc.Queries) gin.HandlerFunc {
	return func(c *gin.Context) {
		followAction(c, false, redisConn, queries)
	}
}
