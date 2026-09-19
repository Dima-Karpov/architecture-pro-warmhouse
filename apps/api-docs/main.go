package main

import "github.com/gofiber/fiber/v3"

//	@title			Тёплый дом — REST API
//	@version		1.0.0
//	@description	Синхронный HTTP JSON через API Gateway (ADR-006).
//	@servers.url	https://api.warmhouse.example
//	@servers.description	API Gateway

func main() {
	app := fiber.New()
	app.Get("/api/v1/devices/:id", GetDevice)
	app.Patch("/api/v1/devices/:id/status", UpdateDeviceStatus)
	app.Post("/api/v1/commands", SendCommandHandler)
	app.Get("/api/v1/telemetry", GetTelemetry)
	_ = app.Listen(":8090")
}
