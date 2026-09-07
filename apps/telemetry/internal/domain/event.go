package domain

import (
	"time"

	"github.com/google/uuid"
)

const QueueTelemetryReceived = "telemetry.received"

type ReadingReceived struct {
	RecordedAt time.Time
	Unit       string
	DeviceID   uuid.UUID
	HouseID    uuid.UUID
	Value      float64
}

func NewReceived(reading Reading, houseID uuid.UUID) (ReadingReceived, error) {
	if houseID == uuid.Nil {
		return ReadingReceived{}, ErrInvalidPayload
	}

	return ReadingReceived{
		RecordedAt: reading.RecordedAt,
		Unit:       reading.Unit,
		DeviceID:   reading.DeviceID,
		HouseID:    houseID,
		Value:      reading.Value,
	}, nil
}
