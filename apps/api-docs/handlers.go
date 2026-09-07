package main

import (
	"strings"
	"time"

	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
)

func parseUUID(raw string) (uuid.UUID, bool) {
	id, err := uuid.Parse(raw)
	if err != nil {
		return uuid.Nil, false
	}

	return id, true
}

func newID() uuid.UUID {
	id, err := uuid.NewV7()
	if err != nil {
		return uuid.New()
	}

	return id
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

func notFound(id uuid.UUID) NotFoundResponse {
	return NotFoundResponse{Errors: []NotFoundError{{
		Type:    "NotFoundError",
		Message: "Device not found",
		Context: ErrorContext{ID: id},
	}}}
}

func sampleDevice(id uuid.UUID, status string) Device {
	return Device{
		ID:           id,
		TypeID:       parseMust("01932c4e-8a1c-7111-8222-2b5c6d7e8f90"),
		HouseID:      parseMust("01932c4e-8a1d-7222-8333-3c6d7e8f9012"),
		TypeCode:     "heating",
		SerialNumber: "RL-22",
		Address:      "192.168.10.4:47808",
		Status:       status,
	}
}

func parseMust(raw string) uuid.UUID {
	id, _ := uuid.Parse(raw)
	return id
}

// GetDevice godoc
//
//	@Summary		Получить устройство
//	@Description	Жилец смотрит прибор. command тем же методом берёт type_code и address.
//	@Tags			device
//	@Produce		json
//	@Param			id	path		string	true	"Device ID"	format(uuid)
//	@Success		200	{object}	Device
//	@Failure		404	{object}	NotFoundResponse
//	@Failure		500	{object}	InternalErrorResponse
//	@Router			/api/v1/devices/{id} [get]
func GetDevice(c fiber.Ctx) error {
	id, ok := parseUUID(c.Params("id"))
	if !ok {
		return c.Status(fiber.StatusNotFound).JSON(notFound(uuid.Nil))
	}

	return c.JSON(sampleDevice(id, "on"))
}

// UpdateDeviceStatus godoc
//
//	@Summary	Обновить состояние устройства
//	@Tags		device
//	@Accept		json
//	@Produce	json
//	@Param		id		path		string				true	"Device ID"	format(uuid)
//	@Param		request	body		DeviceStatusUpdate	true	"Новый статус on/off"
//	@Success	200		{object}	Device
//	@Failure	404		{object}	NotFoundResponse
//	@Failure	500		{object}	InternalErrorResponse
//	@Router		/api/v1/devices/{id}/status [patch]
func UpdateDeviceStatus(c fiber.Ctx) error {
	id, ok := parseUUID(c.Params("id"))
	if !ok {
		return c.Status(fiber.StatusNotFound).JSON(notFound(uuid.Nil))
	}

	var req DeviceStatusUpdate
	if err := c.Bind().JSON(&req); err != nil {
		return err
	}

	return c.JSON(sampleDevice(id, req.Status))
}

// SendCommand godoc
//
//	@Summary		Отправить команду устройству
//	@Description	Ручная команда или вызов сценария. source ставит сервис (токен жильца или scenario), не клиент. В брокер не идёт.
//	@Tags			command
//	@Accept			json
//	@Produce		json
//	@Param			request	body		SendCommand	true	"Команда"
//	@Success		200		{object}	Command
//	@Failure		404		{object}	NotFoundResponse
//	@Failure		409		{object}	ConflictResponse
//	@Failure		500		{object}	InternalErrorResponse
//	@Router			/api/v1/commands [post]
func SendCommandHandler(c fiber.Ctx) error {
	var req SendCommand
	if err := c.Bind().JSON(&req); err != nil {
		return err
	}

	return c.JSON(Command{
		ID:        newID(),
		DeviceID:  req.DeviceID,
		Action:    req.Action,
		Source:    "user",
		Result:    "ok",
		CreatedAt: time.Now().UTC(),
	})
}

// GetTelemetry godoc
//
//	@Summary	Получить показания устройств
//	@Tags		telemetry
//	@Produce	json
//	@Param		device_ids	query		string	true	"ID устройств через запятую"	example(01932c4e-8a1b-7f3c-8d2e-1a4b5c6d7e8f,01932c4e-8a1c-7111-8222-2b5c6d7e8f90)
//	@Success	200			{array}		TelemetryData
//	@Failure	404			{object}	NotFoundResponse
//	@Failure	500			{object}	InternalErrorResponse
//	@Router		/api/v1/telemetry [get]
func GetTelemetry(c fiber.Ctx) error {
	ids := parseQueryUUIDs(c.Query("device_ids"))
	now := time.Now().UTC()
	out := make([]TelemetryData, 0, len(ids))

	for _, id := range ids {
		out = append(out, TelemetryData{
			ID:         newID(),
			DeviceID:   id,
			Value:      21.4,
			Unit:       "°C",
			RecordedAt: now,
		})
	}

	return c.JSON(out)
}
