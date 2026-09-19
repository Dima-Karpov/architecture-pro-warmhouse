package api

import "github.com/gofiber/fiber/v3"

func Register(app *fiber.App, controller *DeviceController) {
	app.Get("/api/v1/devices", controller.List)
	app.Post("/api/v1/devices", controller.Upsert)
	app.Get("/api/v1/devices/:id", controller.Get)
	app.Patch("/api/v1/devices/:id/status", controller.UpdateStatus)
}
