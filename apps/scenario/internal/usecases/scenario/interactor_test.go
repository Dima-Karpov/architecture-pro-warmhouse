package scenario

import (
	"testing"
	"time"

	"github.com/google/uuid"

	"scenario/internal/domain"
)

func TestHandleTurnsOnWhenCold(t *testing.T) {
	t.Parallel()

	got := NewInteractor().Handle(domain.ReadingReceived{
		RecordedAt: time.Now().UTC(),
		Unit:       "°C",
		DeviceID:   uuid.MustParse("01932c4e-8a1b-7f3c-8d2e-1a4b5c6d7e8f"),
		HouseID:    uuid.MustParse("01932c4e-8a1d-7222-8333-3c6d7e8f9012"),
		Value:      domain.ColdCelsius - 1,
	})

	if got.Action != domain.ActionOn {
		t.Fatalf("action: got %s, want %s", got.Action, domain.ActionOn)
	}
}
