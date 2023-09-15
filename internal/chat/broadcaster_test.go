package chat

import (
	"context"
	"log"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/sesaquecruz/go-chat-broadcaster/config"
	"github.com/sesaquecruz/go-chat-broadcaster/internal/model"
	"github.com/sesaquecruz/go-chat-broadcaster/internal/rabbitmq"
	"github.com/sesaquecruz/go-chat-broadcaster/internal/redis"
	"github.com/sesaquecruz/go-chat-broadcaster/test"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func SetupContainers(ctx context.Context) (*test.RabbitMQContainer, *test.RedisContainer) {
	configDir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	for i := 0; i < 2; i++ {
		configDir = filepath.Dir(configDir)
	}

	rbt := test.NewRabbitMQContainer(ctx, configDir)
	rds := test.NewRedisContainer(ctx)

	return rbt, rds
}

func SetupBrokers(rabbitContainer *test.RabbitMQContainer, redisContainer *test.RedisContainer) (*rabbitmq.Broker, *redis.Broker) {
	cfg := &config.Config{
		RabbitMqUrl: rabbitContainer.Url(),
		RedisAddr:   redisContainer.Addr(),
	}

	conn, ch, err := rabbitmq.Connection(cfg)
	if err != nil {
		log.Fatal(err)
	}

	rdb := redis.Connection(cfg)

	rabbitBroker := rabbitmq.NewBroker(conn, ch)
	redisBroker := redis.NewBroker(rdb)

	return rabbitBroker, redisBroker
}

func TestShouldReceiveAndDeliverMessages(t *testing.T) {
	ctx := context.Background()
	rabbitBroker, redisBroker := SetupBrokers(SetupContainers(ctx))

	err := Broadcast(ctx, rabbitBroker, redisBroker)
	assert.Nil(t, err)

	roomId := uuid.NewString()
	msg := &model.Message{Id: uuid.NewString(), RoomId: roomId}

	go func() {
		for {
			err := rabbitBroker.Publish(ctx, msg)
			assert.Nil(t, err)
			<-time.After(500 * time.Millisecond)
		}
	}()

	subs := 10
	recv := make(chan *model.Message)

	sub := func() {
		msgs := redisBroker.Subscribe(ctx, roomId)
		m := <-msgs
		recv <- m
	}

	for i := 0; i < subs; i++ {
		go sub()
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
