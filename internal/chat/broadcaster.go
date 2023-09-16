package chat

import (
	"context"
	"errors"
	"sync"

	"github.com/sesaquecruz/go-chat-broadcaster/internal/rabbitmq"
	"github.com/sesaquecruz/go-chat-broadcaster/internal/redis"
	"github.com/sesaquecruz/go-chat-broadcaster/pkg/log"
)

var ErrBroadcasterAlreadyStarted = errors.New("broadcaster already started")
var ErrBroadcasterStopped = errors.New("broadcaster stopped")

type Broadcaster struct {
	mu           sync.Mutex
	rabbitBroker *rabbitmq.Broker
	redisBroker  *redis.Broker
	started      bool
	logger       *log.Logger
}

func NewBroadcaster(rabbitBroker *rabbitmq.Broker, redisBroker *redis.Broker) *Broadcaster {
	return &Broadcaster{
		rabbitBroker: rabbitBroker,
		redisBroker:  redisBroker,
		started:      false,
		logger:       log.NewLoggerOfObject(Broadcaster{}),
	}
}

func (b *Broadcaster) Start(ctx context.Context) error {
	b.mu.Lock()

	if b.started {
		b.mu.Unlock()
		err := ErrBroadcasterAlreadyStarted
		b.logger.Error(err)
		return err
	}

	b.started = true
	defer func() { b.started = false }()
	b.mu.Unlock()

	msgs, err := b.rabbitBroker.Subscribe(ctx)
	if err != nil {
		b.logger.Error(err)
		return err
	}

	for msg := range msgs {
		err = b.redisBroker.Publish(ctx, msg)
		if err != nil {
			b.logger.Error(err)
		}
	}

	return ErrBroadcasterStopped
}
