package interfaces

import (
	"strconv"
	"time"

	"example.com/golang_twitter/config"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
)

type RedisConn struct {
	Client *redis.Client
}

func NewRedisConn(cfg config.Config) *RedisConn {
	client := redis.NewClient(&redis.Options{
		Addr:     cfg.RedisAddr,
		Password: cfg.RedisPassword,
		DB:       cfg.RedisDB,
	})
	return &RedisConn{Client: client}
}

func (r *RedisConn) Close() error {
	return r.Client.Close()
}

func (r *RedisConn) SetSession(c *gin.Context, sessionID string, userId int64, expiration time.Duration) error {
	_, err := r.Client.Set(c.Request.Context(), sessionID, userId, expiration).Result()
	return err
}

func (r *RedisConn) GetUserIdFromSession(c *gin.Context, sessionID string) (int64, error) {
	result, err := r.Client.Get(c.Request.Context(), sessionID).Result()
	if err != nil {
		return 0, err
	}

	userId, err := strconv.ParseInt(result, 10, 64)
	if err != nil {
		return 0, err
	}

	return userId, nil
}
