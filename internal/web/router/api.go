package router

import (
	"github.com/sesaquecruz/go-chat-broadcaster/config"
	"github.com/sesaquecruz/go-chat-broadcaster/internal/web/handler"
	"github.com/sesaquecruz/go-chat-broadcaster/pkg/middleware"

	"github.com/gin-gonic/gin"
)

func ApiRouter(
	cfg *config.Config,
	healthz *handler.Healthz,
	subscriber *handler.Subscriber,
) *gin.Engine {

	gin.SetMode("release")

	r := gin.New()
	r.Use(middleware.CorsMiddleware())

	api := r.Group(cfg.ApiPath)
	{
		api.GET("/healthz", gin.WrapH(healthz.Healthz()))

		api.Use(middleware.JwtMiddleware(cfg.JwtIssuer, cfg.JwtAudience))

		api.GET("/subscribe/:roomId", subscriber.Subscribe)
	}

	return r
}
