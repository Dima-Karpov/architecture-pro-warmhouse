package api

import (
	"errors"

	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"

	"device/internal/domain"
	"device/internal/interfaces/api/responses"
	usecase "device/internal/usecases/device"
)

type DeviceController struct {
	usecase DeviceUsecase
}

func NewDeviceController(usecase DeviceUsecase) *DeviceController {
	return &DeviceController{usecase: usecase}
}

// Get godoc
//
//	@Summary		Получить устройство
//	@Description	Карточка прибора. Монолит и command берут type_code и address.
//	@Tags			device
//	@Produce		json
//	@Param			id	path		string	true	"Device ID"	format(uuid)
//	@Success		200	{object}	responses.DeviceResponse
//	@Failure		404	{object}	responses.NotFoundResponse
//	@Failure		500	{object}	responses.InternalErrorResponse
//	@Router			/api/v1/devices/{id} [get]
func (controller *DeviceController) Get(ctx fiber.Ctx) error {
	id, ok := parseUUID(ctx.Params("id"))
	if !ok {
		return ctx.Status(fiber.StatusNotFound).JSON(responses.NotFound(uuid.Nil))
	}

	item, err := controller.usecase.Get(ctx.Context(), id)
	if err != nil {
		return writeUsecaseError(ctx, id, err)
	}

	return ctx.JSON(responses.NewDeviceResponse(item))
}

// List godoc
//
//	@Summary	Список устройств
//	@Tags		device
//	@Produce	json
//	@Success	200	{array}	responses.DeviceResponse
//	@Failure	500	{object}	responses.InternalErrorResponse
//	@Router		/api/v1/devices [get]
func (controller *DeviceController) List(ctx fiber.Ctx) error {
	items, err := controller.usecase.List(ctx.Context())
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(responses.Internal())
	}

	out := make([]responses.DeviceResponse, 0, len(items))
	for _, item := range items {
		out = append(out, responses.NewDeviceResponse(item))
	}

	return ctx.JSON(out)
}

// Upsert godoc
//
//	@Summary		Создать или обновить устройство
//	@Description	Монолит отдаёт реестр sensors: тот же external_id обновляет карточку.
//	@Tags			device
//	@Accept			json
//	@Produce		json
//	@Param			request	body		responses.UpsertRequest	true	"Устройство"
//	@Success		200		{object}	responses.DeviceResponse
//	@Failure		500		{object}	responses.InternalErrorResponse
//	@Router			/api/v1/devices [post]
func (controller *DeviceController) Upsert(ctx fiber.Ctx) error {
	var req responses.UpsertRequest
	if err := ctx.Bind().JSON(&req); err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(responses.Internal())
	}

	item, err := controller.usecase.Upsert(ctx.Context(), usecase.UpsertInput{
		SerialNumber: req.SerialNumber,
		Address:      req.Address,
		TypeCode:     req.TypeCode,
		Status:       req.Status,
		ExternalID:   req.ExternalID,
	})
	if err != nil {
		return writeUsecaseError(ctx, uuid.Nil, err)
	}

	return ctx.JSON(responses.NewDeviceResponse(item))
}

// UpdateStatus godoc
//
//	@Summary	Обновить состояние устройства
//	@Tags		device
//	@Accept		json
//	@Produce	json
//	@Param		id		path		string					true	"Device ID"	format(uuid)
//	@Param		request	body		responses.StatusUpdate	true	"Новый статус on/off"
//	@Success	200		{object}	responses.DeviceResponse
//	@Failure	404		{object}	responses.NotFoundResponse
//	@Failure	500		{object}	responses.InternalErrorResponse
//	@Router		/api/v1/devices/{id}/status [patch]
func (controller *DeviceController) UpdateStatus(ctx fiber.Ctx) error {
	id, ok := parseUUID(ctx.Params("id"))
	if !ok {
		return ctx.Status(fiber.StatusNotFound).JSON(responses.NotFound(uuid.Nil))
	}

	var req responses.StatusUpdate
	if err := ctx.Bind().JSON(&req); err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(responses.Internal())
	}

	item, err := controller.usecase.UpdateStatus(ctx.Context(), id, req.Status)
	if err != nil {
		return writeUsecaseError(ctx, id, err)
	}

	return ctx.JSON(responses.NewDeviceResponse(item))
}

func parseUUID(raw string) (uuid.UUID, bool) {
	id, err := uuid.Parse(raw)
	if err != nil {
		return uuid.Nil, false
	}

	return id, true
}

func writeUsecaseError(ctx fiber.Ctx, id uuid.UUID, err error) error {
	if errors.Is(err, domain.ErrNotFound) {
		return ctx.Status(fiber.StatusNotFound).JSON(responses.NotFound(id))
	}

	return ctx.Status(fiber.StatusInternalServerError).JSON(responses.Internal())
}
