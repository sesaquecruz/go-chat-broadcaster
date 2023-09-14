package redis

import "github.com/redis/go-redis/v9"

func Connection(addr string) *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: "",
		DB:       0,
	})
}
