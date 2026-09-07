package main

import (
	"context"
	"log"

	"github.com/gofiber/contrib/v3/swaggerui"
	"github.com/gofiber/fiber/v3"

	"telemetry/docs"
	"telemetry/internal/infrastructure/amqp"
	"telemetry/internal/infrastructure/clock"
	"telemetry/internal/infrastructure/config"
	"telemetry/internal/infrastructure/db"
	"telemetry/internal/infrastructure/healthcheck"
	"telemetry/internal/infrastructure/idgen"
	"telemetry/internal/interfaces/api"
	usecase "telemetry/internal/usecases/telemetry"
)

//	@title			telemetry
//	@version		1.0.0
//	@description	Показания устройств. Монолит отдаёт опрос температуры на переход.
//	@servers.url	http://localhost:8083
//	@servers.description	local

func main() {
	cfg := config.Load()
	ctx := context.Background()

	conn, err := db.Open(ctx, cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("database: %v", err)
	}

	publisher, err := amqp.Open(cfg.AMQPURL)
	if err != nil {
		db.Close(conn)
		log.Fatalf("broker: %v", err)
	}

	defer db.Close(conn)
	defer publisher.Close()

	interactor := usecase.NewInteractor(
		db.NewStore(conn),
		idgen.NewV7(),
		clock.NewUTC(),
		publisher,
		cfg.HouseID,
	)
	controller := api.NewTelemetryController(interactor)

	app := fiber.New()
	healthcheck.Register(app)
	api.Register(app, controller)
	app.Use(swaggerui.New(swaggerui.Config{
		BasePath:    "/",
		Path:        "swagger",
		FileContent: docs.Spec,
		Title:       "telemetry",
	}))

	if err = app.Listen(cfg.Addr); err != nil {
		log.Printf("listen: %v", err)
	}
}
