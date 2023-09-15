package redis

import (
	"github.com/sesaquecruz/go-chat-broadcaster/config"

	"github.com/redis/go-redis/v9"
)

func Connection(cfg *config.Config) *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr:     cfg.RedisAddr,
		Password: "",
		DB:       0,
	})
}
