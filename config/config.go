package config

import (
	"strings"

	"github.com/spf13/viper"
)

type Config struct {
	ServiceName    string
	ServiceVersion string
	ApiPath        string
	ApiPort        string
	JwtIssuer      string
	JwtAudience    []string
	RabbitMqUrl    string
	RedisUrl       string
}

var (
	env  *viper.Viper
	file *viper.Viper
	cfg  *Config
)

func init() {
	env = viper.New()
	env.SetDefault("APP_VERSION", "")
	env.SetDefault("APP_API_PORT", "")
	env.SetDefault("APP_JWT_ISSUER", "")
	env.SetDefault("APP_JWT_AUDIENCE", "")
	env.SetDefault("APP_RABBITMQ_URL", "")
	env.SetDefault("APP_REDIS_URL", "")
	env.AutomaticEnv()

	file = viper.New()
	file.SetConfigName("config")
	file.SetConfigType("toml")
	file.AddConfigPath(".")
	file.ReadInConfig()
}

func getEnv(key string) string {
	envVal := env.GetString(key)
	if envVal != "" {
		return envVal
	}

	fileVal := file.GetString(strings.Replace(key, "_", ".", -1))
	return fileVal
}

func Load() *Config {
	if cfg == nil {
		cfg = &Config{
			ServiceName:    "go-chat-broadcaster",
			ServiceVersion: getEnv("APP_VERSION"),
			ApiPath:        "/api/v1",
			ApiPort:        getEnv("APP_API_PORT"),
			JwtIssuer:      getEnv("APP_JWT_ISSUER"),
			JwtAudience:    strings.Split(getEnv("APP_JWT_AUDIENCE"), ","),
			RabbitMqUrl:    getEnv("APP_RABBITMQ_URL"),
			RedisUrl:       getEnv("APP_REDIS_URL"),
		}
	}

	return cfg
}
