package database

import (
	"context"
	"fmt"

	"github.com/codebayu/account-service/internal/config"
	"github.com/redis/go-redis/v9"
)

var Ctx = context.Background()

func NewRedisClient(cfg *config.Config) (*redis.Client, error) {
	addr := fmt.Sprintf("%s:%s", cfg.RedisHost, cfg.RedisPort)
	rdb := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: cfg.RedisPassword,
		DB:       cfg.RedisDB,
	})

	if err := rdb.Ping(Ctx).Err(); err != nil {
		return nil, err
	}

	fmt.Println("✅ redis connected")
	return rdb, nil
}
