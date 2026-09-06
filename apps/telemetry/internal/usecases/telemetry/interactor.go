package telemetry

import (
	"context"
	"time"

	"github.com/google/uuid"

	"telemetry/internal/domain"
)

type Interactor struct {
	repo    Repository
	ids     IDGenerator
	clock   Clock
	events  EventPublisher
	houseID uuid.UUID
}

func NewInteractor(
	repo Repository,
	ids IDGenerator,
	clock Clock,
	events EventPublisher,
	houseID uuid.UUID,
) *Interactor {
	return &Interactor{
		repo:    repo,
		ids:     ids,
		clock:   clock,
		events:  events,
		houseID: houseID,
	}
}

type IngestInput struct {
	RecordedAt time.Time
	Unit       string
	DeviceID   uuid.UUID
	HouseID    uuid.UUID
	Value      float64
}

func (interactor *Interactor) Ingest(ctx context.Context, in IngestInput) (domain.Reading, error) {
	at := in.RecordedAt
	if at.IsZero() {
		at = interactor.clock.Now()
	}

	reading, err := domain.NewReading(interactor.ids.New(), in.DeviceID, in.Value, in.Unit, at)
	if err != nil {
		return domain.Reading{}, err
	}

	saved, err := interactor.repo.Save(ctx, reading)
	if err != nil {
		return domain.Reading{}, err
	}

	houseID := in.HouseID
	if houseID == uuid.Nil {
		houseID = interactor.houseID
	}

	event, err := domain.NewReceived(saved, houseID)
	if err != nil {
		return saved, err
	}

	return saved, interactor.events.PublishReceived(ctx, event)
}

func (interactor *Interactor) List(ctx context.Context, ids []uuid.UUID) ([]domain.Reading, error) {
	if len(ids) == 0 {
		return nil, domain.ErrNotFound
	}

	return interactor.repo.ListByDeviceIDs(ctx, ids)
}
