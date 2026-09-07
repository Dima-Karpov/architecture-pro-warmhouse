package responses

import (
	"time"

	"github.com/google/uuid"

	"telemetry/internal/domain"
)

type ReadingResponse struct {
	RecordedAt time.Time `json:"recorded_at" format:"date-time" example:"2026-09-05T10:00:00Z"`
	Unit       string    `json:"unit" example:"°C"`
	Value      float64   `json:"value" example:"21.4"`
	ID         uuid.UUID `json:"id" swaggertype:"string" example:"01932c4e-8a1f-7444-8555-5e8f90123456"`
	DeviceID   uuid.UUID `json:"device_id" swaggertype:"string" example:"01932c4e-8a1b-7f3c-8d2e-1a4b5c6d7e8f"`
}

type IngestRequest struct {
	RecordedAt time.Time `json:"recorded_at" format:"date-time" example:"2026-09-05T10:00:00Z"`
	Unit       string    `json:"unit" example:"°C"`
	Value      float64   `json:"value" example:"21.4"`
	DeviceID   uuid.UUID `json:"device_id" swaggertype:"string" example:"01932c4e-8a1b-7f3c-8d2e-1a4b5c6d7e8f"`
	HouseID    uuid.UUID `json:"house_id" swaggertype:"string" example:"01932c4e-8a1d-7222-8333-3c6d7e8f9012"`
}

type ErrorContext struct {
	ID uuid.UUID `json:"id" format:"uuid" swaggertype:"string" example:"01932c4e-8a1b-7f3c-8d2e-1a4b5c6d7e8f"`
}

type NotFoundError struct {
	Type    string       `json:"type" example:"NotFoundError" enums:"NotFoundError"`
	Message string       `json:"message" example:"Device not found"`
	Context ErrorContext `json:"context"`
}

type NotFoundResponse struct {
	Errors []NotFoundError `json:"errors"`
}

type InternalServerError struct {
	Type    string `json:"type" example:"InternalServerError" enums:"InternalServerError"`
	Message string `json:"message" example:"Something went wrong"`
}

type InternalErrorResponse struct {
	Errors []InternalServerError `json:"errors"`
}

func NewReadingResponse(item domain.Reading) ReadingResponse {
	return ReadingResponse{
		RecordedAt: item.RecordedAt,
		Unit:       item.Unit,
		Value:      item.Value,
		ID:         item.ID,
		DeviceID:   item.DeviceID,
	}
}

func NotFound(id uuid.UUID) NotFoundResponse {
	return NotFoundResponse{Errors: []NotFoundError{{
		Type:    "NotFoundError",
		Message: "Device not found",
		Context: ErrorContext{ID: id},
	}}}
}

func Internal() InternalErrorResponse {
	return InternalErrorResponse{Errors: []InternalServerError{{
		Type:    "InternalServerError",
		Message: "Something went wrong",
	}}}
}
