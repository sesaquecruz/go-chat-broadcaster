package rabbitmq

import (
	"context"
	"encoding/json"

	"github.com/sesaquecruz/go-chat-broadcaster/internal/model"
	"github.com/sesaquecruz/go-chat-broadcaster/pkg/log"

	amqp "github.com/rabbitmq/amqp091-go"
)

type Broker struct {
	conn   *amqp.Connection
	ch     *amqp.Channel
	logger *log.Logger
}

func NewBroker(conn *amqp.Connection, ch *amqp.Channel) *Broker {
	return &Broker{
		conn:   conn,
		ch:     ch,
		logger: log.NewLoggerOfObject(Broker{}),
	}
}

func (b *Broker) Publish(ctx context.Context, message *model.Message) error {
	body, err := json.Marshal(message)
	if err != nil {
		b.logger.Error(err)
		return err
	}

	msg := amqp.Publishing{
		ContentType: "application/json",
		Body:        body,
	}

	err = b.ch.PublishWithContext(
		ctx,
		"messages",
		"",
		false,
		false,
		msg,
	)
	if err != nil {
		b.logger.Error(err)
		return err
	}

	return nil
}

func (b *Broker) Subscribe(ctx context.Context) (<-chan *model.Message, error) {
	ch, err := b.conn.Channel()
	if err != nil {
		b.logger.Error(err)
		return nil, err
	}

	msgs, err := ch.Consume(
		"messages.queue",
		"",
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		ch.Close()
		b.logger.Error(err)
		return nil, err
	}

	res := make(chan *model.Message)

	go func() {
		defer ch.Close()
		defer close(res)

		for {
			select {
			case <-ctx.Done():
				return
			case msg := <-msgs:
				var message model.Message

				err := json.Unmarshal(msg.Body, &message)
				if err != nil {
					b.logger.Error(err)
				} else {
					res <- &message
				}

				msg.Ack(true)
			}
		}
	}()

	return res, nil
}
