package chat

import (
	"context"
	"log"
	"os"
	"path/filepath"
	"sync"
	"testing"
	"time"

	"github.com/sesaquecruz/go-chat-broadcaster/internal/model"
	"github.com/sesaquecruz/go-chat-broadcaster/internal/rabbitmq"
	"github.com/sesaquecruz/go-chat-broadcaster/internal/redis"
	"github.com/sesaquecruz/go-chat-broadcaster/test"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

const testTimeout = 240 * time.Second

var mu sync.Mutex

var (
	ctx               context.Context
	rabbitMqContainer *test.RabbitMqContainer
	redisContainer    *test.RedisContainer
	rabbitMqBroker    *rabbitmq.Broker
	redisBroker       *redis.Broker
	broadcaster       *Broadcaster
)

func setupBroadcaster() {
	mu.Lock()
	defer mu.Unlock()

	if ctx == nil {
		ctx, _ = context.WithTimeout(context.Background(), testTimeout)
	}

	if rabbitMqContainer == nil {
		currentDepth := 2

		configDir, err := os.Getwd()
		if err != nil {
			log.Fatal(err)
		}

		for i := 0; i < currentDepth; i++ {
			configDir = filepath.Dir(configDir)
		}

		rabbitMqContainer = test.NewRabbitMqContainer(ctx, configDir)
	}

	if redisContainer == nil {
		redisContainer = test.NewRedisContainer(ctx)
	}

	if rabbitMqBroker == nil {
		conn, err := rabbitmq.Connect(rabbitMqContainer.Url())
		if err != nil {
			log.Fatal(err)
		}

		rabbitMqBroker = rabbitmq.NewBroker(conn)
	}

	if redisBroker == nil {
		conn := redis.Connect(redisContainer.Url())
		redisBroker = redis.NewBroker(conn)
	}

	if broadcaster == nil {
		broadcaster = NewBroadcaster(rabbitMqBroker, redisBroker)

		go func() {
			err := broadcaster.Start(ctx)
			if err != nil {
				log.Fatal(err)
			}
		}()

		<-time.After(5 * time.Second)
	}
}

func TestShouldReturnAnErrorWhenCallStartOnStartedBroadcaster(t *testing.T) {
	setupBroadcaster()

	err := broadcaster.Start(ctx)
	assert.ErrorIs(t, err, ErrBroadcasterAlreadyStarted)
}

func TestShouldReceiveAndDeliverMessages(t *testing.T) {
	setupBroadcaster()

	roomId := uuid.NewString()
	msg := &model.Message{Id: uuid.NewString(), RoomId: roomId}

	// Publish on rabbitmq
	go func() {
		for {
			err := rabbitMqBroker.Publish(ctx, msg)
			assert.Nil(t, err)
			<-time.After(500 * time.Millisecond)
		}
	}()

	subs := 10
	recv := make(chan *model.Message)

	// Receive from redis
	for i := 0; i < subs; i++ {
		go func() {
			msgs := redisBroker.Subscribe(ctx, roomId)
			m := <-msgs
			recv <- m
		}()
	}

	timeout := time.After(30 * time.Second)

	for i := 0; i < subs; {
		select {
		case r := <-recv:
			assert.Equal(t, msg.Id, r.Id)
			assert.Equal(t, msg.RoomId, r.RoomId)
			i++
		case <-timeout:
			t.Error("timeout reached")
			return
		}
	}
}
