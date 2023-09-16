package redis

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/sesaquecruz/go-chat-broadcaster/internal/model"
	"github.com/sesaquecruz/go-chat-broadcaster/test"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

const testTimeout = 240 * time.Second

var mu sync.Mutex

var (
	ctx            context.Context
	redisContainer *test.RedisContainer
	broker         *Broker
)

func setupBroker() {
	mu.Lock()
	defer mu.Unlock()

	if ctx == nil {
		ctx, _ = context.WithTimeout(context.Background(), testTimeout)
	}

	if redisContainer == nil {
		redisContainer = test.NewRedisContainer(ctx)
	}

	if broker == nil {
		conn := Connect(redisContainer.Url())
		broker = NewBroker(conn)
	}
}

func TestShouldSendAndReceiveMessages(t *testing.T) {
	setupBroker()

	roomId := uuid.NewString()
	msgs := broker.Subscribe(ctx, roomId)

	for i := 0; i < 10; i++ {
		msg := &model.Message{Id: uuid.NewString(), RoomId: roomId}
		err := broker.Publish(ctx, msg)
		assert.Nil(t, err)

		select {
		case res := <-msgs:
			assert.Equal(t, msg.Id, res.Id)
			assert.Equal(t, msg.RoomId, res.RoomId)
		case <-ctx.Done():
			t.Error("timeout reached")
			return
		}
	}
}
