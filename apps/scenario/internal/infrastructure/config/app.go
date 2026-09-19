package config

import (
	"os"
	"strings"
)

const (
	defaultPort = "8084"
	defaultAMQP = "amqp://guest:guest@localhost:5672/"
)

type App struct {
	AMQPURL string
	Addr    string
}

func Load() App {
	return App{
		AMQPURL: envOr("AMQP_URL", defaultAMQP),
		Addr:    listenAddr(envOr("PORT", defaultPort)),
	}
}

func envOr(key, fallback string) string {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}

	return value
}

func listenAddr(port string) string {
	if !strings.HasPrefix(port, ":") {
		return ":" + port
	}

	return port
}
