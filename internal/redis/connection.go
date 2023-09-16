package redis

import (
	"github.com/redis/go-redis/v9"
)

type Connection struct {
	Rdb *redis.Client
}

func Connect(url string) *Connection {
	rdb := redis.NewClient(&redis.Options{
		Addr:     url,
		Password: "",
		DB:       0,
	})

	return &Connection{rdb}
}
