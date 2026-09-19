package main

import (
	"time"

	"github.com/google/uuid"
)

type Device struct {
	TypeCode     string    `json:"type_code" example:"heating" enums:"heating,light,gate,camera,unknown"`
	SerialNumber string    `json:"serial_number" example:"RL-22"`
	Address      string    `json:"address" example:"192.168.10.4:47808"`
	Status       string    `json:"status" example:"on" enums:"on,off,inactive"`
	ID           uuid.UUID `json:"id" format:"uuid" swaggertype:"string" example:"01932c4e-8a1b-7f3c-8d2e-1a4b5c6d7e8f"`
	TypeID       uuid.UUID `json:"type_id" format:"uuid" swaggertype:"string" example:"01932c4e-8a1c-7111-8222-2b5c6d7e8f90"`
	HouseID      uuid.UUID `json:"house_id" format:"uuid" swaggertype:"string" example:"01932c4e-8a1d-7222-8333-3c6d7e8f9012"`
} // @name Device

type DeviceStatusUpdate struct {
	Status string `json:"status" example:"off" enums:"on,off" validate:"required"`
} // @name DeviceStatusUpdate

type SendCommand struct {
	Action   string    `json:"action" example:"on" enums:"on,off,lock" validate:"required"`
	DeviceID uuid.UUID `json:"device_id" format:"uuid" swaggertype:"string" example:"01932c4e-8a1b-7f3c-8d2e-1a4b5c6d7e8f" validate:"required"`
} // @name SendCommand

type Command struct {
	CreatedAt time.Time `json:"created_at" format:"date-time" example:"2026-09-05T10:00:00Z"`
	Action    string    `json:"action" example:"on" enums:"on,off,lock"`
	Source    string    `json:"source" example:"user" enums:"user,scenario"`
	Result    string    `json:"result" example:"ok"`
	ID        uuid.UUID `json:"id" format:"uuid" swaggertype:"string" example:"01932c4e-8a1e-7333-8444-4d7e8f901234"`
	DeviceID  uuid.UUID `json:"device_id" format:"uuid" swaggertype:"string" example:"01932c4e-8a1b-7f3c-8d2e-1a4b5c6d7e8f"`
} // @name Command

type TelemetryData struct {
	RecordedAt time.Time `json:"recorded_at" format:"date-time" example:"2026-09-05T10:00:00Z"`
	Unit       string    `json:"unit" example:"°C"`
	Value      float64   `json:"value" example:"21.4"`
	ID         uuid.UUID `json:"id" format:"uuid" swaggertype:"string" example:"01932c4e-8a1f-7444-8555-5e8f90123456"`
	DeviceID   uuid.UUID `json:"device_id" format:"uuid" swaggertype:"string" example:"01932c4e-8a1b-7f3c-8d2e-1a4b5c6d7e8f"`
} // @name TelemetryData

type TelemetryReceived struct {
	RecordedAt time.Time `json:"recorded_at" format:"date-time" example:"2026-09-05T10:00:00Z" validate:"required"`
	Unit       string    `json:"unit" example:"°C" validate:"required"`
	Value      float64   `json:"value" example:"17.2" validate:"required"`
	DeviceID   uuid.UUID `json:"device_id" validate:"required" format:"uuid" swaggertype:"string" example:"01932c4e-8a1b-7f3c-8d2e-1a4b5c6d7e8f"`
	HouseID    uuid.UUID `json:"house_id" validate:"required" format:"uuid" swaggertype:"string" example:"01932c4e-8a1d-7222-8333-3c6d7e8f9012"`
} // @name TelemetryReceived

type ErrorContext struct {
	ID uuid.UUID `json:"id" format:"uuid" swaggertype:"string" example:"01932c4e-8a1b-7f3c-8d2e-1a4b5c6d7e8f"`
} // @name ErrorContext

type NotFoundError struct {
	Type    string       `json:"type" example:"NotFoundError" enums:"NotFoundError"`
	Message string       `json:"message" example:"Device not found"`
	Context ErrorContext `json:"context"`
} // @name NotFoundError

type NotFoundResponse struct {
	Errors []NotFoundError `json:"errors"`
} // @name NotFoundResponse

type DeviceInactiveError struct {
	Type    string       `json:"type" example:"DeviceInactiveError" enums:"DeviceInactiveError"`
	Message string       `json:"message" example:"Device is inactive"`
	Context ErrorContext `json:"context"`
} // @name DeviceInactiveError

type ConflictResponse struct {
	Errors []DeviceInactiveError `json:"errors"`
} // @name ConflictResponse

type InternalServerError struct {
	Type    string `json:"type" example:"InternalServerError" enums:"InternalServerError"`
	Message string `json:"message" example:"Something went wrong"`
} // @name InternalServerError

type InternalErrorResponse struct {
	Errors []InternalServerError `json:"errors"`
} // @name InternalErrorResponse
