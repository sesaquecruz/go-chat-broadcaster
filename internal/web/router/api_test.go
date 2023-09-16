package router

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/sesaquecruz/go-chat-broadcaster/config"
	"github.com/sesaquecruz/go-chat-broadcaster/internal/chat"
	"github.com/sesaquecruz/go-chat-broadcaster/internal/model"
	"github.com/sesaquecruz/go-chat-broadcaster/internal/rabbitmq"
	"github.com/sesaquecruz/go-chat-broadcaster/internal/redis"
	"github.com/sesaquecruz/go-chat-broadcaster/internal/web/handler"
	"github.com/sesaquecruz/go-chat-broadcaster/test"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

const testTimeout = 240 * time.Second

var mu sync.Mutex

var (
	ctx               context.Context
	cfg               *config.Config
	rabbitMqContainer *test.RabbitMqContainer
	redisContainer    *test.RedisContainer
	authServer        *test.AuthServer
	rabbitMqConn      *rabbitmq.Connection
	redisConn         *redis.Connection
	rabbitMqBroker    *rabbitmq.Broker
	redisBroker       *redis.Broker
	broadcaster       *chat.Broadcaster
	healthzHandler    *handler.Healthz
	subscriberHandler *handler.Subscriber
	apiRouter         *gin.Engine
	apiUrl            string
)

func setupApiRouter() {
	mu.Lock()
	defer mu.Unlock()

	if ctx == nil {
		ctx, _ = context.WithTimeout(context.Background(), testTimeout)
	}

	if cfg == nil {
		cfg = &config.Config{
			ServiceName:    "test-name",
			ServiceVersion: "test-version",
		}
	}

	if rabbitMqContainer == nil {
		currentDepth := 3

		configDir, err := os.Getwd()
		if err != nil {
			log.Fatal(err)
		}

		for i := 0; i < currentDepth; i++ {
			configDir = filepath.Dir(configDir)
		}

		rabbitMqContainer = test.NewRabbitMqContainer(ctx, configDir)
		cfg.RabbitMqUrl = rabbitMqContainer.Url()
	}

	if redisContainer == nil {
		redisContainer = test.NewRedisContainer(ctx)
		cfg.RedisUrl = redisContainer.Url()
	}

	if authServer == nil {
		authServer = test.NewAuth0Server()
		cfg.JwtIssuer = authServer.GetIssuer()
		cfg.JwtAudience = authServer.GetAudience()
	}

	if broadcaster == nil {
		var err error

		rabbitMqConn, err = rabbitmq.Connect(cfg.RabbitMqUrl)
		if err != nil {
			log.Fatal(err)
		}

		redisConn = redis.Connect(cfg.RedisUrl)

		rabbitMqBroker = rabbitmq.NewBroker(rabbitMqConn)
		redisBroker = redis.NewBroker(redisConn)

		broadcaster = chat.NewBroadcaster(rabbitMqBroker, redisBroker)

		go func() {
			err = broadcaster.Start(ctx)
			if err != nil {
				log.Fatal(err)
			}
		}()
	}

	if healthzHandler == nil {
		healthzHandler = handler.NewHealthz(cfg, rabbitMqConn, redisConn)
	}

	if subscriberHandler == nil {
		subscriberHandler = handler.NewSubscriber(redisBroker)
	}

	if apiRouter == nil {
		listener, err := net.Listen("tcp", ":0")
		if err != nil {
			log.Fatal(err)
		}

		port := listener.Addr().(*net.TCPAddr).Port
		if err := listener.Close(); err != nil {
			log.Fatal(err)
		}

		cfg.ApiPath = "/api/v1"
		apiUrl = fmt.Sprintf("http://127.0.0.1:%d%s", port, cfg.ApiPath)

		apiRouter = ApiRouter(cfg, healthzHandler, subscriberHandler)

		go func() {
			err := apiRouter.Run(fmt.Sprintf(":%d", port))
			if err != nil {
				log.Fatal(err)
			}
		}()

		<-time.After(10 * time.Second)
	}
}

func TestShouldReturnHealthInfo(t *testing.T) {
	setupApiRouter()

	type Component struct {
		Name    string `json:"name"`
		Version string `json:"version"`
	}

	type Info struct {
		Status    string    `json:"status"`
		Timestamp string    `json:"timestamp"`
		Component Component `json:"component"`
	}

	req, err := http.NewRequest(http.MethodGet, apiUrl+"/healthz", nil)
	assert.Nil(t, err)

	client := http.Client{}
	res, err := client.Do(req)
	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, res.StatusCode)

	var info Info
	err = json.NewDecoder(res.Body).Decode(&info)
	assert.Nil(t, err)
	res.Body.Close()

	assert.Equal(t, cfg.ServiceName, info.Component.Name)
	assert.Equal(t, cfg.ServiceVersion, info.Component.Version)
}

func TestShouldReturnUnauthorizedWhenTrySubscribeWithoutAuthorization(t *testing.T) {
	setupApiRouter()

	url := fmt.Sprintf("%s/subscribe/%s", apiUrl, uuid.NewString())
	req, err := http.NewRequest(http.MethodGet, url, nil)
	assert.Nil(t, err)

	client := http.Client{}
	res, err := client.Do(req)
	assert.Nil(t, err)
	assert.Equal(t, http.StatusUnauthorized, res.StatusCode)
}

func TestShouldSubscribeAndReceiveMessages(t *testing.T) {
	setupApiRouter()

	roomId := uuid.NewString()

	// Publish messages with sequential ids
	go func() {
		i := 0

		for {
			msg := &model.Message{Id: strconv.Itoa(i), RoomId: roomId}
			err := rabbitMqBroker.Publish(ctx, msg)
			assert.Nil(t, err)
			<-time.After(1000 * time.Millisecond)

			i++
		}
	}()

	// Subscribe on the room
	type Event struct {
		Event string        `json:"event"`
		Data  model.Message `json:"data"`
	}

	parseEvent := func(s string) (e Event, err error) {
		msgs := strings.Split(s, "\n")
		event := strings.Replace(msgs[0], "event:", "", 1)
		data := strings.Replace(msgs[1], "data:", "", 1)
		result := fmt.Sprintf(`{"event": "%s", "data": %s}`, event, data)
		err = json.Unmarshal([]byte(result), &e)
		return
	}

	url := fmt.Sprintf("%s/subscribe/%s", apiUrl, roomId)

	jwt, err := authServer.GenerateJwt(authServer.GenerateSubject())
	assert.Nil(t, err)

	req, err := http.NewRequest(http.MethodGet, url, nil)
	assert.Nil(t, err)

	req.Header.Set("Authorization", "Bearer "+jwt)

	client := http.Client{}
	res, err := client.Do(req)
	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, res.StatusCode)

	// Success on receive 10 sequential messages
	buffer := make([]byte, 256)

	first := true
	total := 0
	previousId := 0

	for total < 10 {
		select {
		case <-ctx.Done():
			t.Error("timeout reached")
			return
		default:
			n, err := res.Body.Read(buffer)
			assert.Nil(t, err)

			event, err := parseEvent(string(buffer[:n]))
			assert.Nil(t, err)

			assert.Equal(t, "message", event.Event)
			assert.Equal(t, roomId, event.Data.RoomId)

			msgId, err := strconv.Atoi(event.Data.Id)
			assert.Nil(t, err)

			if first {
				first = false
				total = 1
			} else {
				if msgId-1 == previousId {
					total++
				} else {
					total = 0
				}
			}

			previousId = msgId
		}
	}
}
