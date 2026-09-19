package main

import (
	"log"

	"github.com/gofiber/contrib/v3/swaggerui"
	"github.com/gofiber/fiber/v3"

	"temperature-api/docs"
	"temperature-api/internal/infrastructure/clock"
	"temperature-api/internal/infrastructure/config"
	"temperature-api/internal/infrastructure/healthcheck"
	"temperature-api/internal/infrastructure/random"
	"temperature-api/internal/interfaces/api"
	"temperature-api/internal/usecases/temperature"
)

//	@title			temperature-api
//	@version		1.0.0
//	@description	Имитация датчика температуры для монолита smart_home.
//	@servers.url	http://localhost:8081
//	@servers.description	local

func main() {
	cfg := config.Load()
	interactor := temperature.NewInteractor(clock.NewUTC(), random.NewSource())
	controller := api.NewTemperatureController(interactor)

	app := fiber.New()
	healthcheck.Register(app)
	api.Register(app, controller)
	app.Use(swaggerui.New(swaggerui.Config{
		BasePath:    "/",
		Path:        "swagger",
		FileContent: docs.Spec,
		Title:       "temperature-api",
	}))

	log.Fatal(app.Listen(cfg.Addr))
}
