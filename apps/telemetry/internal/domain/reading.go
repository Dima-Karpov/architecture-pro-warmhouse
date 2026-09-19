package domain

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

const UnitCelsius = "°C"

var (
	ErrNotFound       = errors.New("device not found")
	ErrInvalidPayload = errors.New("invalid payload")
)

type Reading struct {
	RecordedAt time.Time
	Unit       string
	ID         uuid.UUID
	DeviceID   uuid.UUID
	Value      float64
}

func NewReading(id, deviceID uuid.UUID, value float64, unit string, at time.Time) (Reading, error) {
	if deviceID == uuid.Nil {
		return Reading{}, ErrInvalidPayload
	}

	if unit == "" {
		unit = UnitCelsius
	}

	if at.IsZero() {
		return Reading{}, ErrInvalidPayload
	}

	return Reading{
		RecordedAt: at.UTC(),
		Unit:       unit,
		ID:         id,
		DeviceID:   deviceID,
		Value:      value,
	}, nil
}
