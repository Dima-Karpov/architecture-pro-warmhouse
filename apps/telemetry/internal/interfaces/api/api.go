package api

import "github.com/gofiber/fiber/v3"

func Register(app *fiber.App, controller *TelemetryController) {
	app.Get("/api/v1/telemetry", controller.List)
	app.Post("/api/v1/telemetry", controller.Ingest)
}
