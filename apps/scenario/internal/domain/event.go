package domain

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

const (
	QueueTelemetryReceived = "telemetry.received"

	ActionOn  = "on"
	ActionOff = "off"

	ColdCelsius = 18.0
	HotCelsius  = 25.0
)

var ErrInvalidPayload = errors.New("invalid payload")

type ReadingReceived struct {
	RecordedAt time.Time
	Unit       string
	DeviceID   uuid.UUID
	HouseID    uuid.UUID
	Value      float64
}

type Decision struct {
	Action string
	Reason string
}

func ParseReceived(
	deviceID, houseID uuid.UUID,
	value float64,
	unit string,
	at time.Time,
) (ReadingReceived, error) {
	if deviceID == uuid.Nil || houseID == uuid.Nil || at.IsZero() {
		return ReadingReceived{}, ErrInvalidPayload
	}

	return ReadingReceived{
		RecordedAt: at.UTC(),
		Unit:       unit,
		DeviceID:   deviceID,
		HouseID:    houseID,
		Value:      value,
	}, nil
}

func Decide(value float64) Decision {
	if value < ColdCelsius {
		return Decision{Action: ActionOn, Reason: "cold"}
	}

	if value > HotCelsius {
		return Decision{Action: ActionOff, Reason: "hot"}
	}

	return Decision{Reason: "in range"}
}
