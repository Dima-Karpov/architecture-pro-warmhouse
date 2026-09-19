package config

import (
	"os"
	"strings"
)

const defaultPort = "8081"

type App struct {
	Addr string
}

func Load() App {
	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	if !strings.HasPrefix(port, ":") {
		port = ":" + port
	}

	return App{Addr: port}
}
