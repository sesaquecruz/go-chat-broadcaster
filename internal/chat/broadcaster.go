package chat

import (
	"context"

	"github.com/sesaquecruz/go-chat-broadcaster/internal/rabbitmq"
	"github.com/sesaquecruz/go-chat-broadcaster/internal/redis"
	"github.com/sesaquecruz/go-chat-broadcaster/pkg/log"
)

func Broadcast(ctx context.Context, rabbitmq *rabbitmq.Broker, redis *redis.Broker) error {
	logger := log.NewLogger("chat.Broadcast")

	msgs, err := rabbitmq.Subscribe(ctx)
	if err != nil {
		return err
	}

	go func() {
		for msg := range msgs {
			err = redis.Publish(ctx, msg)
			if err != nil {
				logger.Error(err)
			}
		}
	}()

	return nil
}
