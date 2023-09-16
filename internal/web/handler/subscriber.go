package handler

import (
	"encoding/json"
	"io"

	"github.com/sesaquecruz/go-chat-broadcaster/internal/chat"
	"github.com/sesaquecruz/go-chat-broadcaster/pkg/log"

	"github.com/gin-gonic/gin"
)

type Subscriber struct {
	redisSub chat.Subscriber
	logger   *log.Logger
}

func NewSubscriber(redisSub chat.Subscriber) *Subscriber {
	return &Subscriber{
		redisSub: redisSub,
		logger:   log.NewLoggerOfObject(Subscriber{}),
	}
}

func (h *Subscriber) Subscribe(c *gin.Context) {
	msgs := h.redisSub.Subscribe(c.Request.Context(), c.Param("roomId"))

	c.Writer.Header().Set("Content-Type", "text/event-stream")
	c.Writer.Header().Set("Cache-Control", "no-cache")
	c.Writer.Header().Set("Transfer-Encoding", "chunked")
	c.Writer.Header().Set("Connection", "keep-alive")

	c.Stream(func(w io.Writer) bool {
		select {
		case msg := <-msgs:
			data, err := json.Marshal(msg)
			if err != nil {
				h.logger.Error(err)
			} else {
				c.SSEvent("message", string(data))
			}
			return true
		default:
			return true
		}
	})
}
