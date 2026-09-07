package api

import (
	"github.com/gofiber/fiber/v3"

	"temperature-api/internal/interfaces/api/responses"
)

type TemperatureController struct {
	getter TemperatureGetter
}

func NewTemperatureController(getter TemperatureGetter) *TemperatureController {
	return &TemperatureController{getter: getter}
}

// ByLocation godoc
//
//	@Summary		Температура по комнате
//	@Description	Случайное показание. Пустой location мапится по sensor id: 1 Living Room, 2 Bedroom, 3 Kitchen.
//	@Tags			temperature
//	@Produce		json
//	@Param			location	query		string	false	"Living Room, Bedroom, Kitchen"
//	@Success		200			{object}	responses.TemperatureResponse
//	@Router			/temperature [get]
func (controller *TemperatureController) ByLocation(ctx fiber.Ctx) error {
	return ctx.JSON(responses.NewTemperatureResponse(
		controller.getter.Get(ctx.Query("location"), ""),
	))
}

// ByID godoc
//
//	@Summary		Температура по id датчика
//	@Description	1 Living Room, 2 Bedroom, 3 Kitchen. Иначе location=Unknown.
//	@Tags			temperature
//	@Produce		json
//	@Param			id	path		string	true	"sensor id"
//	@Success		200	{object}	responses.TemperatureResponse
//	@Router			/temperature/{id} [get]
func (controller *TemperatureController) ByID(ctx fiber.Ctx) error {
	return ctx.JSON(responses.NewTemperatureResponse(
		controller.getter.Get("", ctx.Params("id")),
	))
}
