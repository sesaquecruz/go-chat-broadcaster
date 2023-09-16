package rabbitmq

import (
	"context"
	"log"
	"os"
	"path/filepath"
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
	ctx               context.Context
	rabbitMqContainer *test.RabbitMqContainer
	broker            *Broker
)

func setupBroker() {
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

	if broker == nil {
		conn, err := Connect(rabbitMqContainer.Url())
		if err != nil {
			log.Fatal(err)
		}

		broker = NewBroker(conn)
	}
}

func TestShouldSendAndReceiveMessages(t *testing.T) {
	setupBroker()

	msgs, err := broker.Subscribe(ctx)
	assert.Nil(t, err)

	for i := 0; i < 10; i++ {
		msg := &model.Message{Id: uuid.NewString(), RoomId: uuid.NewString()}
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
