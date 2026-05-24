package config

import (
	"os"
	"strings"

	"github.com/joho/godotenv"
)

type Config struct {
	Port           string
	DBPath         string
	RedisURL       string
	RabbitMQURL    string
	AllowedOrigins []string
	Environment    string
}

func Load() *Config {
	_ = godotenv.Load() // Ignore error if .env doesn't exist

	allowedOrigins := os.Getenv("ALLOWED_ORIGINS")
	origins := []string{"*"}
	if allowedOrigins != "" {
		origins = strings.Split(allowedOrigins, ",")
	}

	return &Config{
		Port:           getEnv("PORT", "8080"),
		DBPath:         getEnv("DB_PATH", "construction.db"),
		RedisURL:       getEnv("REDIS_URL", "localhost:6379"),
		RabbitMQURL:    getEnv("RABBITMQ_URL", "amqp://guest:guest@localhost:5672/"),
		AllowedOrigins: origins,
		Environment:    getEnv("ENVIRONMENT", "development"),
	}
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}
