package chat

import (
	"context"

	"github.com/sesaquecruz/go-chat-broadcaster/internal/model"
)

type Subscriber interface {
	Subscribe(ctx context.Context, roomId string) <-chan *model.Message
}
