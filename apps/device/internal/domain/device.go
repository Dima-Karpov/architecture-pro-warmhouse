package domain

import (
	"errors"

	"github.com/google/uuid"
)

const (
	StatusOn       = "on"
	StatusOff      = "off"
	StatusInactive = "inactive"

	TypeHeating = "heating"
	TypeLight   = "light"
	TypeGate    = "gate"
	TypeCamera  = "camera"
	TypeUnknown = "unknown"
)

var (
	ErrNotFound       = errors.New("device not found")
	ErrInvalidStatus  = errors.New("invalid status")
	ErrInvalidType    = errors.New("invalid type_code")
	ErrInvalidPayload = errors.New("invalid payload")
)

type Device struct {
	SerialNumber string
	Address      string
	TypeCode     string
	Status       string
	ExternalID   string
	ID           uuid.UUID
	TypeID       uuid.UUID
	HouseID      uuid.UUID
}

func ParseStatus(raw string) (string, error) {
	switch raw {
	case StatusOn, StatusOff, StatusInactive:
		return raw, nil
	default:
		return "", ErrInvalidStatus
	}
}

func ParseTypeCode(raw string) (string, error) {
	switch raw {
	case TypeHeating, TypeLight, TypeGate, TypeCamera, TypeUnknown:
		return raw, nil
	default:
		return "", ErrInvalidType
	}
}
