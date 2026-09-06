package responses

import (
	"github.com/google/uuid"

	"device/internal/domain"
)

type DeviceResponse struct {
	TypeCode     string    `json:"type_code" example:"heating" enums:"heating,light,gate,camera,unknown"`
	SerialNumber string    `json:"serial_number" example:"RL-22"`
	Address      string    `json:"address" example:"192.168.10.4:47808"`
	Status       string    `json:"status" example:"on" enums:"on,off,inactive"`
	ExternalID   string    `json:"external_id,omitempty" example:"1"`
	ID           uuid.UUID `json:"id" swaggertype:"string" example:"01932c4e-8a1b-7f3c-8d2e-1a4b5c6d7e8f"`
	TypeID       uuid.UUID `json:"type_id" swaggertype:"string" example:"01932c4e-8a1c-7111-8222-2b5c6d7e8f90"`
	HouseID      uuid.UUID `json:"house_id" swaggertype:"string" example:"01932c4e-8a1d-7222-8333-3c6d7e8f9012"`
}

type StatusUpdate struct {
	Status string `json:"status" example:"off" enums:"on,off"`
}

type UpsertRequest struct {
	SerialNumber string `json:"serial_number" example:"RL-22"`
	Address      string `json:"address" example:"Living Room"`
	TypeCode     string `json:"type_code" example:"heating" enums:"heating,light,gate,camera,unknown"`
	Status       string `json:"status" example:"on" enums:"on,off,inactive"`
	ExternalID   string `json:"external_id" example:"1"`
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

func NewDeviceResponse(item domain.Device) DeviceResponse {
	return DeviceResponse{
		TypeCode:     item.TypeCode,
		SerialNumber: item.SerialNumber,
		Address:      item.Address,
		Status:       item.Status,
		ExternalID:   item.ExternalID,
		ID:           item.ID,
		TypeID:       item.TypeID,
		HouseID:      item.HouseID,
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
