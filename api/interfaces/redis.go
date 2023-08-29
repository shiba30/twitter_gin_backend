package interfaces

import (
	"context"
	"time"

	"example.com/golang_twitter/config"
	"github.com/go-redis/redis/v8"
)

var ctx = context.Background()

type RedisConn struct {
	Client *redis.Client
}

func InitializeRedis(cfg config.Config) *RedisConn {
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

func (r *RedisConn) SetSession(sessionID string, email string, expiration time.Duration) error {
	_, err := r.Client.Set(ctx, sessionID, email, expiration).Result()
	return err
}

func (r *RedisConn) GetSession(sessionID string) (string, error) {
	return r.Client.Get(ctx, sessionID).Result()
}
