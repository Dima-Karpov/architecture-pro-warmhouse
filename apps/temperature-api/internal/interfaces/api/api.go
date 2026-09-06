package api

import "github.com/gofiber/fiber/v3"

func Register(app *fiber.App, controller *TemperatureController) {
	app.Get("/temperature", controller.ByLocation)
	app.Get("/temperature/:id", controller.ByID)
}
