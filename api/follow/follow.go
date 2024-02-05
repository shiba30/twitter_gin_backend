package follow

import (
	"log"

	"example.com/golang_twitter/config"
	sqlc "example.com/golang_twitter/db/sqlc"
	"example.com/golang_twitter/utils"
	"github.com/gin-gonic/gin"
)

type followForm struct {
	UserId int64 `json:"userId"`
}

func followAction(c *gin.Context, isFollow bool, cfg config.Config, queries *sqlc.Queries) {
	form := followForm{}
	if err := c.ShouldBind(&form); err != nil {
		log.Printf("failed to bind data: %v", err)
		c.JSON(400, gin.H{"error": "処理に失敗しました"})
		return
	}

	sessionID, err := c.Cookie("session_id")
	if err != nil {
		log.Printf("Failed to retrieve session ID: %v", err)
		c.JSON(500, gin.H{"error": "セッションIDの取得に失敗しました"})
		return
	}

	userID, err := utils.GetSessionUserId(c, sessionID)
	if err != nil {
		log.Printf("Failed to retrieve userId from session: %v", err)
		c.JSON(500, gin.H{"error": "セッションからユーザー情報の取得に失敗しました"})
		return
	}

	if isFollow {
		// フォロー
		err = queries.CreateFollow(c, sqlc.CreateFollowParams{
			FollowerID: userID,
			FolloweeID: form.UserId,
		})
	} else {
		// フォロー解除
		err = queries.DeleteFollow(c, sqlc.DeleteFollowParams{
			FollowerID: userID,
			FolloweeID: form.UserId,
		})
	}

	if err != nil {
		log.Printf("Error in follow action: %v", err)
		c.JSON(500, gin.H{"error": "処理に失敗しました"})
		return
	}

	log.Printf("successful follow execute")
	c.JSON(200, gin.H{"message": "フォロー処理に成功しました"})
}

// フォロー機能
func Follow(cfg config.Config, queries *sqlc.Queries) gin.HandlerFunc {
	return func(c *gin.Context) {
		followAction(c, true, cfg, queries)
	}
}

// フォロー解除機能
func UnFollow(cfg config.Config, queries *sqlc.Queries) gin.HandlerFunc {
	return func(c *gin.Context) {
		followAction(c, false, cfg, queries)
	}
}
