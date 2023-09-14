package test

import (
	"context"
	"fmt"
	"log"

	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/redis"
)

type RedisContainer struct {
	container *redis.RedisContainer
	addr      string
}

func NewRedisContainer(ctx context.Context) *RedisContainer {
	container, err := redis.RunContainer(ctx, testcontainers.WithImage("redis:7.2-alpine"))
	if err != nil {
		log.Fatal(err)
	}

	host, err := container.Host(ctx)
	if err != nil {
		log.Fatal(err)
	}

	port, err := container.MappedPort(ctx, "6379/tcp")
	if err != nil {
		log.Fatal(err)
	}

	return &RedisContainer{
		container: container,
		addr:      fmt.Sprintf("%s:%s", host, port.Port()),
	}
}

func (c *RedisContainer) Addr() string {
	return c.addr
}
