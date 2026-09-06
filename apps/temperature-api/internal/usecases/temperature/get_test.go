package temperature

import (
	"testing"
	"time"

	"temperature-api/internal/domain"
)

type stubClock struct {
	now time.Time
}

func (s stubClock) Now() time.Time {
	return s.now
}

type stubValues struct {
	value float64
}

func (s stubValues) Next() float64 {
	return s.value
}

func TestInteractorGetMapsLivingRoom(t *testing.T) {
	t.Parallel()

	at := time.Date(2026, 9, 6, 12, 0, 0, 0, time.UTC)
	interactor := NewInteractor(stubClock{now: at}, stubValues{value: 21.5})

	got := interactor.Get("Living Room", "")

	if got.SensorID != domain.IDLivingRoom {
		t.Fatalf("sensor id: got %s, want %s", got.SensorID, domain.IDLivingRoom)
	}

	if got.Value != 21.5 {
		t.Fatalf("value: got %v, want 21.5", got.Value)
	}

	if !got.Timestamp.Equal(at) {
		t.Fatalf("timestamp: got %s, want %s", got.Timestamp, at)
	}
}
