package config

import (
	"log"
	"os"
)

type Config struct {
	ServiceName string
	ApiVersion  string
	ApiPort     string
	JwtIssuer   string
	JwtAudience string
	RabbitMqUrl string
	RedisAddr   string
}

var cfg *Config

func getEnv(key string) string {
	value, ok := os.LookupEnv(key)
	if !ok {
		log.Fatalf("env var %s is required\n", key)
	}

	return value
}

func init() {
	cfg = &Config{
		ServiceName: "go-chat-broadaster",
		ApiVersion:  "v1",
		ApiPort:     getEnv("BROAD_API_PORT"),
		JwtIssuer:   getEnv("BROAD_JWT_ISSUER"),
		JwtAudience: getEnv("BROAD_JWT_AUDIENCE"),
		RabbitMqUrl: getEnv("BROAD_RABBITMQ_URL"),
		RedisAddr:   getEnv("BROAD_REDIS_ADDR"),
	}
}

func GetConfig() *Config {
	return cfg
}
