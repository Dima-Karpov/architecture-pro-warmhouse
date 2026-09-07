package main

import (
	"context"
	"log"

	"github.com/gofiber/contrib/v3/swaggerui"
	"github.com/gofiber/fiber/v3"

	"device/docs"
	"device/internal/infrastructure/config"
	"device/internal/infrastructure/db"
	"device/internal/infrastructure/healthcheck"
	"device/internal/infrastructure/idgen"
	"device/internal/interfaces/api"
	usecase "device/internal/usecases/device"
)

//	@title			device
//	@version		1.0.0
//	@description	Реестр устройств. Монолит отдаёт sensors на переход.
//	@servers.url	http://localhost:8082
//	@servers.description	local

func main() {
	cfg := config.Load()
	ctx := context.Background()

	conn, err := db.Open(ctx, cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("database: %v", err)
	}
	defer db.Close(conn)

	interactor := usecase.NewInteractor(db.NewStore(conn), idgen.NewV7(), cfg.HouseID, cfg.TypeID)
	controller := api.NewDeviceController(interactor)

	app := fiber.New()
	healthcheck.Register(app)
	api.Register(app, controller)
	app.Use(swaggerui.New(swaggerui.Config{
		BasePath:    "/",
		Path:        "swagger",
		FileContent: docs.Spec,
		Title:       "device",
	}))

	if err = app.Listen(cfg.Addr); err != nil {
		log.Printf("listen: %v", err)
	}
}
