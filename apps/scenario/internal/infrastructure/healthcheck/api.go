package healthcheck

import "github.com/gofiber/fiber/v3"

type Response struct {
	Status string `json:"status"`
}

func Register(app *fiber.App) {
	app.Get("/health", health)
}

func health(ctx fiber.Ctx) error {
	return ctx.JSON(Response{Status: "ok"})
}
