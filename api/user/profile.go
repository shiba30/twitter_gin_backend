package user

import (
	"database/sql"
	"log"
	"strconv"
	"time"

	"example.com/golang_twitter/api/interfaces"
	"example.com/golang_twitter/config"
	sqlc "example.com/golang_twitter/db/sqlc"
	"example.com/golang_twitter/utils"
	"github.com/gin-gonic/gin"
)

type profileForm struct {
	HeaderImage  string `json:"header_image"`
	ProfileImage string `json:"profile_image"`
	DisplayName  string `json:"display_name" binding:"required"`
	Bio          string `json:"bio"`
	Location     string `json:"location"`
	Website      string `json:"website"`
	BirthDay     string `json:"birthday"`
}

func ProfileRoutes(router *gin.RouterGroup, cfg config.Config, redisConn *interfaces.RedisConn, queries *sqlc.Queries) {
	router.GET("/profile/:id", GetProfile(redisConn, queries))
	router.POST("/profile", SaveProfile(cfg, redisConn, queries))
}

func GetProfile(redisConn *interfaces.RedisConn, queries *sqlc.Queries) gin.HandlerFunc {
	return func(c *gin.Context) {
		// ログインユーザ情報の取得
		userInfo, err := utils.CurrentUser(c, redisConn, queries)
		if err != nil {
			c.JSON(500, gin.H{"error": "ログインユーザ情報の取得に失敗しました"})
			return
		}

		// ユーザIDの取得
		var userId int64
		userIdStr := c.Param("id")
		if userIdStr != "null" {
			userId, err = strconv.ParseInt(userIdStr, 10, 64)
			if err != nil {
				c.JSON(500, gin.H{"error": "ユーザIDの取得に失敗しました"})
				return
			}
		} else {
			userId = userInfo.ID
		}

		// プロフィール情報の取得
		profile, err := queries.GetUserInfo(c, userId)
		if err != nil {
			c.JSON(500, gin.H{"error": "プロフィール情報の取得に失敗しました"})
			return
		}

		// ツイート情報の取得
		tweets, err := queries.GetTweetsByUser(c, userId)
		if err != nil {
			log.Printf("Failed to retrieve tweets: %v", err)
			c.JSON(500, gin.H{"error": "ツイートの取得に失敗しました"})
			return
		}

		c.JSON(200, gin.H{
			"profile":       profile,
			"tweets":        tweets,
			"currentUserId": userInfo.ID,
		})
	}
}

func SaveProfile(cfg config.Config, redisConn *interfaces.RedisConn, queries *sqlc.Queries) gin.HandlerFunc {
	return func(c *gin.Context) {
		form := profileForm{}

		if err := c.ShouldBindJSON(&form); err != nil {
			// バリデーションエラー時の処理
			c.JSON(400, gin.H{"error": "入力値が正しくありません"})
			return
		}

		// ユーザ情報の取得
		userInfo, err := utils.CurrentUser(c, redisConn, queries)
		if err != nil {
			c.JSON(500, gin.H{"error": "プロフィール情報の取得に失敗しました"})
			return
		}

		bio := toNullString(form.Bio)
		location := toNullString(form.Location)
		website := toNullString(form.Website)
		birthDate, err := toNullTime(form.BirthDay)
		if err != nil {
			handleError(c, 400, "誕生日の形式が正しくありません")
			return
		}

		// 画像処理部分を共通関数に置き換え
		headerImagePath, err := utils.ProcessImage(form.HeaderImage, cfg.UploadedImagesDir, userInfo.ID)
		if err != nil {
			c.JSON(500, gin.H{"error": "画像の処理に失敗しました"})
			return
		}
		profileImagePath, err := utils.ProcessImage(form.ProfileImage, cfg.UploadedImagesDir, userInfo.ID)
		if err != nil {
			c.JSON(500, gin.H{"error": "画像の処理に失敗しました"})
			return
		}

		// プロフィール情報の更新
		profile, err := queries.UpdateUserProfile(c, sqlc.UpdateUserProfileParams{
			ID:           userInfo.ID,
			DisplayName:  form.DisplayName,
			Bio:          bio,
			Location:     location,
			Website:      website,
			BirthDate:    birthDate,
			HeaderImage:  headerImagePath,
			ProfileImage: profileImagePath,
		})
		if err != nil {
			// データベースからの取得に失敗した場合のエラー処理
			c.JSON(500, gin.H{"error": "プロフィール情報の更新に失敗しました"})
			return
		}
		c.JSON(200, gin.H{"profile": profile})
	}
}

func handleError(c *gin.Context, statusCode int, message string) {
	log.Println(message)
	c.JSON(statusCode, gin.H{"error": message})
}

func toNullString(s string) sql.NullString {
	if s == "" {
		return sql.NullString{Valid: false}
	}
	return sql.NullString{String: s, Valid: true}
}

func toNullTime(dateStr string) (sql.NullTime, error) {
	if dateStr == "" {
		return sql.NullTime{Valid: false}, nil
	}
	parsedDate, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		return sql.NullTime{}, err
	}
	return sql.NullTime{Time: parsedDate, Valid: true}, nil
}
