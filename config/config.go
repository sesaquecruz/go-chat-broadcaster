package config

import (
	"log"
	"os"
	"strings"
)

type Info struct {
	ServiceName string
	ApiVersion  string
}

type Config struct {
	ApiPort     string
	JwtIssuer   string
	JwtAudience []string
	RabbitMqUrl string
	RedisAddr   string
}

var info *Info
var cfg *Config

func getEnv(key string) string {
	value, ok := os.LookupEnv(key)
	if !ok {
		log.Printf("env var %s was not found\n", key)
	}

	return value
}

func init() {
	info = &Info{
		ServiceName: "go-chat-broadaster",
		ApiVersion:  "v1",
	}

	cfg = &Config{
		ApiPort:     getEnv("BROAD_API_PORT"),
		JwtIssuer:   getEnv("BROAD_JWT_ISSUER"),
		JwtAudience: strings.Split(getEnv("BROAD_JWT_AUDIENCE"), ","),
		RabbitMqUrl: getEnv("BROAD_RABBITMQ_URL"),
		RedisAddr:   getEnv("BROAD_REDIS_ADDR"),
	}
}

func GetInfo() *Info {
	return info
}

func GetConfig() *Config {
	return cfg
}
