package config

import (
	"os"
	"strings"

	"github.com/google/uuid"
)

const (
	defaultPort    = "8082"
	defaultDB      = "postgres://postgres:postgres@localhost:5432/device"
	defaultHouseID = "01932c4e-8a1d-7222-8333-3c6d7e8f9012"
	defaultTypeID  = "01932c4e-8a1c-7111-8222-2b5c6d7e8f90"
)

type App struct {
	DatabaseURL string
	Addr        string
	HouseID     uuid.UUID
	TypeID      uuid.UUID
}

func Load() App {
	return App{
		DatabaseURL: envOr("DATABASE_URL", defaultDB),
		Addr:        listenAddr(envOr("PORT", defaultPort)),
		HouseID:     parseUUID(envOr("DEFAULT_HOUSE_ID", defaultHouseID)),
		TypeID:      parseUUID(envOr("DEFAULT_TYPE_ID", defaultTypeID)),
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
