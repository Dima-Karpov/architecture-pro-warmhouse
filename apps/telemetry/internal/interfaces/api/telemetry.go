package api

import (
	"errors"
	"strings"

	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"

	"telemetry/internal/domain"
	"telemetry/internal/interfaces/api/responses"
	usecase "telemetry/internal/usecases/telemetry"
)

type TelemetryController struct {
	usecase TelemetryUsecase
}

func NewTelemetryController(usecase TelemetryUsecase) *TelemetryController {
	return &TelemetryController{usecase: usecase}
}

// List godoc
//
//	@Summary	Получить показания устройств
//	@Tags		telemetry
//	@Produce	json
//	@Param		device_ids	query		string	true	"ID устройств через запятую"	example(01932c4e-8a1b-7f3c-8d2e-1a4b5c6d7e8f)
//	@Success	200			{array}		responses.ReadingResponse
//	@Failure	404			{object}	responses.NotFoundResponse
//	@Failure	500			{object}	responses.InternalErrorResponse
//	@Router		/api/v1/telemetry [get]
func (controller *TelemetryController) List(ctx fiber.Ctx) error {
	ids := parseQueryUUIDs(ctx.Query("device_ids"))
	if len(ids) == 0 {
		return ctx.Status(fiber.StatusNotFound).JSON(responses.NotFound(uuid.Nil))
	}

	items, err := controller.usecase.List(ctx.Context(), ids)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			return ctx.Status(fiber.StatusNotFound).JSON(responses.NotFound(uuid.Nil))
		}

		return ctx.Status(fiber.StatusInternalServerError).JSON(responses.Internal())
	}

	out := make([]responses.ReadingResponse, 0, len(items))
	for _, item := range items {
		out = append(out, responses.NewReadingResponse(item))
	}

	return ctx.JSON(out)
}

// Ingest godoc
//
//	@Summary		Записать показание
//	@Description	Монолит отдаёт опрос temperature-api на переход.
//	@Tags			telemetry
//	@Accept			json
//	@Produce		json
//	@Param			request	body		responses.IngestRequest	true	"Показание"
//	@Success		200		{object}	responses.ReadingResponse
//	@Failure		500		{object}	responses.InternalErrorResponse
//	@Router			/api/v1/telemetry [post]
func (controller *TelemetryController) Ingest(ctx fiber.Ctx) error {
	var req responses.IngestRequest
	if err := ctx.Bind().JSON(&req); err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(responses.Internal())
	}

	item, err := controller.usecase.Ingest(ctx.Context(), usecase.IngestInput{
		RecordedAt: req.RecordedAt,
		Unit:       req.Unit,
		DeviceID:   req.DeviceID,
		HouseID:    req.HouseID,
		Value:      req.Value,
	})
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(responses.Internal())
	}

	return ctx.JSON(responses.NewReadingResponse(item))
}

func parseQueryUUIDs(raw string) []uuid.UUID {
	if raw == "" {
		return nil
	}

	parts := strings.Split(raw, ",")
	ids := make([]uuid.UUID, 0, len(parts))

	for _, part := range parts {
		id, err := uuid.Parse(strings.TrimSpace(part))
		if err != nil {
			continue
		}

		ids = append(ids, id)
	}

	return ids
}
