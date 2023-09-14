package rabbitmq

import (
	"context"
	"log"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/sesaquecruz/go-chat-broadcaster/internal/model"
	"github.com/sesaquecruz/go-chat-broadcaster/test"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func SetupContainer(ctx context.Context) *test.RabbitMQContainer {
	configDir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	for i := 0; i < 2; i++ {
		configDir = filepath.Dir(configDir)
	}

	container := test.NewRabbitMQContainer(ctx, configDir)
	return container
}

func SetupBroker(container *test.RabbitMQContainer) *Broker {
	conn, ch, err := Connection(container.Url())
	if err != nil {
		log.Fatal(err)
	}

	broker := NewBroker(conn, ch)
	return broker
}

func TestShouldSendAndReceiveMessages(t *testing.T) {
	ctx := context.Background()
	broker := SetupBroker(SetupContainer(ctx))

	msgs, err := broker.Subscribe(ctx)
	assert.Nil(t, err)

	timeout := time.After(30 * time.Second)

	for i := 0; i < 10; i++ {
		msg := &model.Message{Id: uuid.NewString(), RoomId: uuid.NewString()}
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
