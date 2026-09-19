package config

import (
	"os"
	"strings"

	"github.com/google/uuid"
)

const (
	defaultPort    = "8083"
	defaultDB      = "postgres://postgres:postgres@localhost:5432/telemetry"
	defaultAMQP    = "amqp://guest:guest@localhost:5672/"
	defaultHouseID = "01932c4e-8a1d-7222-8333-3c6d7e8f9012"
)

type App struct {
	DatabaseURL string
	AMQPURL     string
	Addr        string
	HouseID     uuid.UUID
}

func Load() App {
	return App{
		DatabaseURL: envOr("DATABASE_URL", defaultDB),
		AMQPURL:     envOr("AMQP_URL", defaultAMQP),
		Addr:        listenAddr(envOr("PORT", defaultPort)),
		HouseID:     parseUUID(envOr("DEFAULT_HOUSE_ID", defaultHouseID)),
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

func parseUUID(raw string) uuid.UUID {
	id, err := uuid.Parse(raw)
	if err != nil {
		return uuid.Nil
	}

	return id
}
