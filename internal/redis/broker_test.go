package redis

import (
	"context"
	"testing"
	"time"

	"github.com/sesaquecruz/go-chat-broadcaster/config"
	"github.com/sesaquecruz/go-chat-broadcaster/internal/model"
	"github.com/sesaquecruz/go-chat-broadcaster/test"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func SetupContainer(ctx context.Context) *test.RedisContainer {
	container := test.NewRedisContainer(ctx)
	return container
}

func SetupBroker(container *test.RedisContainer) *Broker {
	rdb := Connection(&config.Config{RedisAddr: container.Addr()})
	broker := NewBroker(rdb)
	return broker
}

func TestShouldSendAndReceiveMessages(t *testing.T) {
	ctx := context.Background()
	broker := SetupBroker(SetupContainer(ctx))

	roomId := uuid.NewString()
	msgs := broker.Subscribe(ctx, roomId)

	timeout := time.After(30 * time.Second)

	for i := 0; i < 10; i++ {
		msg := &model.Message{Id: uuid.NewString(), RoomId: roomId}
		err := broker.Publish(ctx, msg)
		assert.Nil(t, err)

		select {
		case res := <-msgs:
			assert.Equal(t, msg.Id, res.Id)
			assert.Equal(t, msg.RoomId, res.RoomId)
		case <-timeout:
			t.Error("timeout reached")
			return
		}
	}
}
