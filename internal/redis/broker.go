package redis

import (
	"context"
	"encoding/json"

	"github.com/sesaquecruz/go-chat-broadcaster/internal/model"
	"github.com/sesaquecruz/go-chat-broadcaster/pkg/log"
)

type Broker struct {
	conn   *Connection
	logger *log.Logger
}

func NewBroker(conn *Connection) *Broker {
	return &Broker{
		conn:   conn,
		logger: log.NewLoggerOfObject(Broker{}),
	}
}

func (b *Broker) Publish(ctx context.Context, message *model.Message) error {
	data, err := json.Marshal(message)
	if err != nil {
		b.logger.Error(err)
		return err
	}

	err = b.conn.Rdb.Publish(ctx, message.RoomId, data).Err()
	if err != nil {
		b.logger.Error(err)
		return err
	}

	return nil
}

func (b *Broker) Subscribe(ctx context.Context, roomId string) <-chan *model.Message {
	sub := b.conn.Rdb.Subscribe(ctx, roomId)
	msgs := sub.Channel()
	res := make(chan *model.Message)

	go func() {
		defer sub.Close()
		defer close(res)

		for {
			select {
			case <-ctx.Done():
				return
			case msg := <-msgs:
				var message model.Message

				err := json.Unmarshal([]byte(msg.Payload), &message)
				if err != nil {
					b.logger.Error(err)
				} else {
					res <- &message
				}
			}
		}
	}()

	return res
}
