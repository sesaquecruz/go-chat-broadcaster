package main

import (
	"context"
	"fmt"

	"github.com/sesaquecruz/go-chat-broadcaster/config"
	"github.com/sesaquecruz/go-chat-broadcaster/di"
	"github.com/sesaquecruz/go-chat-broadcaster/internal/rabbitmq"
	"github.com/sesaquecruz/go-chat-broadcaster/internal/redis"
	"github.com/sesaquecruz/go-chat-broadcaster/pkg/log"
)

func main() {
	logger := log.NewLogger("main")
	cfg := config.Load()

	rabbitMqConn, err := rabbitmq.Connect(cfg.RabbitMqUrl)
	if err != nil {
		logger.Fatal(err)
	}

	redisConn := redis.Connect(cfg.RedisUrl)

	chatBroadcaster := di.NewChatBroadcaster(rabbitMqConn, redisConn)
	apiRouter := di.NewApiRouter(cfg, rabbitMqConn, redisConn)

	go func() {
		err = chatBroadcaster.Start(context.Background())
		if err != nil {
			logger.Fatal(err)
		}
	}()

	logger.Infof("api path: %s\n", cfg.ApiPath)
	logger.Infof("running on port %s\n", cfg.ApiPort)

	apiRouter.Run(fmt.Sprintf(":%s", cfg.ApiPort))
}
