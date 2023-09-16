//go:build wireinject
// +build wireinject

package di

import (
	"github.com/sesaquecruz/go-chat-broadcaster/config"
	"github.com/sesaquecruz/go-chat-broadcaster/internal/chat"
	"github.com/sesaquecruz/go-chat-broadcaster/internal/rabbitmq"
	"github.com/sesaquecruz/go-chat-broadcaster/internal/redis"
	"github.com/sesaquecruz/go-chat-broadcaster/internal/web/handler"
	"github.com/sesaquecruz/go-chat-broadcaster/internal/web/router"

	"github.com/gin-gonic/gin"
	"github.com/google/wire"
)

var setSubscriber = wire.NewSet(
	redis.NewBroker,
	wire.Bind(new(chat.Subscriber), new(*redis.Broker)),
)

func NewChatBroadcaster(
	rabbitMqConn *rabbitmq.Connection,
	redisConn *redis.Connection,
) (b *chat.Broadcaster) {

	wire.Build(
		rabbitmq.NewBroker,
		redis.NewBroker,

		chat.NewBroadcaster,
	)

	return
}

func NewApiRouter(
	cfg *config.Config,
	rabbitMqConn *rabbitmq.Connection,
	redisConn *redis.Connection,
) (r *gin.Engine) {

	wire.Build(
		setSubscriber,

		handler.NewHealthz,
		handler.NewSubscriber,

		router.ApiRouter,
	)

	return
}
