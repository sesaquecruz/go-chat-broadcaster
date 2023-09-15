package rabbitmq

import (
	"github.com/sesaquecruz/go-chat-broadcaster/config"

	amqp "github.com/rabbitmq/amqp091-go"
)

func Connection(cfg *config.Config) (*amqp.Connection, *amqp.Channel, error) {
	conn, err := amqp.Dial(cfg.RabbitMqUrl)
	if err != nil {
		return nil, nil, err
	}

	ch, err := conn.Channel()
	if err != nil {
		return nil, nil, err
	}

	return conn, ch, nil
}
