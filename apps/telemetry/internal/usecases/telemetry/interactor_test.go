package telemetry

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"

	"telemetry/internal/domain"
)

type stubRepo struct {
	items []domain.Reading
}

func (s *stubRepo) Save(_ context.Context, reading domain.Reading) (domain.Reading, error) {
	s.items = append(s.items, reading)

	return reading, nil
}

func (s *stubRepo) ListByDeviceIDs(_ context.Context, ids []uuid.UUID) ([]domain.Reading, error) {
	want := make(map[uuid.UUID]struct{}, len(ids))
	for _, id := range ids {
		want[id] = struct{}{}
	}

	out := make([]domain.Reading, 0)
	for _, item := range s.items {
		if _, ok := want[item.DeviceID]; ok {
			out = append(out, item)
		}
	}

	return out, nil
}

type stubIDs struct {
	id uuid.UUID
}

func (s stubIDs) New() uuid.UUID {
	return s.id
}

type stubClock struct {
	now time.Time
}

func (s stubClock) Now() time.Time {
	return s.now
}

type stubPublisher struct {
	events []domain.ReadingReceived
}

func (s *stubPublisher) PublishReceived(_ context.Context, event domain.ReadingReceived) error {
	s.events = append(s.events, event)

	return nil
}

func TestInteractorIngestAndList(t *testing.T) {
	t.Parallel()

	deviceID := uuid.MustParse("01932c4e-8a1b-7f3c-8d2e-1a4b5c6d7e8f")
	houseID := uuid.MustParse("01932c4e-8a1d-7222-8333-3c6d7e8f9012")
	readingID := uuid.MustParse("01932c4e-8a1f-7444-8555-5e8f90123456")
	at := time.Date(2026, 9, 6, 12, 0, 0, 0, time.UTC)
	repo := &stubRepo{}
	publisher := &stubPublisher{}
	interactor := NewInteractor(repo, stubIDs{id: readingID}, stubClock{now: at}, publisher, houseID)

	got, err := interactor.Ingest(context.Background(), IngestInput{
		DeviceID: deviceID,
		Value:    21.4,
		Unit:     domain.UnitCelsius,
	})
	if err != nil {
		t.Fatalf("ingest: %v", err)
	}

	if got.ID != readingID || got.Value != 21.4 || !got.RecordedAt.Equal(at) {
		t.Fatalf("unexpected reading: %+v", got)
	}

	listed, err := interactor.List(context.Background(), []uuid.UUID{deviceID})
	if err != nil {
		t.Fatalf("list: %v", err)
	}

	if len(listed) != 1 {
		t.Fatalf("list len: got %d", len(listed))
	}

	if len(publisher.events) != 1 {
		t.Fatalf("published: got %d", len(publisher.events))
	}

	if publisher.events[0].HouseID != houseID || publisher.events[0].DeviceID != deviceID {
		t.Fatalf("event: %+v", publisher.events[0])
	}
}
