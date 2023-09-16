package handler

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/sesaquecruz/go-chat-broadcaster/config"
	"github.com/sesaquecruz/go-chat-broadcaster/internal/rabbitmq"
	"github.com/sesaquecruz/go-chat-broadcaster/internal/redis"

	"github.com/hellofresh/health-go/v5"
)

const timeout = 5 * time.Second

type Healthz struct {
	health *health.Health
}

func NewHealthz(
	cfg *config.Config,
	rabbitmConn *rabbitmq.Connection,
	radisClient *redis.Connection,
) *Healthz {

	h, _ := health.New(
		health.WithComponent(health.Component{
			Name:    cfg.ServiceName,
			Version: cfg.ServiceVersion,
		}),
	)

	h.Register(
		health.Config{
			Name:    "rabbitmq",
			Timeout: timeout,
			Check: func(ctx context.Context) error {
				if rabbitmConn.Conn.IsClosed() {
					return errors.New("closed connection")
				}
				return nil
			},
		},
	)

	h.Register(
		health.Config{
			Name:    "redis",
			Timeout: timeout,
			Check: func(ctx context.Context) error {
				_, err := radisClient.Rdb.Ping(ctx).Result()
				return err
			},
		},
	)

	return &Healthz{
		health: h,
	}
}

func (h *Healthz) Healthz() http.Handler {
	return h.health.Handler()
}
