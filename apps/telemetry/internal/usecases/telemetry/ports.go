package telemetry

import (
	"context"
	"time"

	"github.com/google/uuid"

	"telemetry/internal/domain"
)

type Repository interface {
	Save(ctx context.Context, reading domain.Reading) (domain.Reading, error)
	ListByDeviceIDs(ctx context.Context, ids []uuid.UUID) ([]domain.Reading, error)
}

type IDGenerator interface {
	New() uuid.UUID
}

type Clock interface {
	Now() time.Time
}

type EventPublisher interface {
	PublishReceived(ctx context.Context, event domain.ReadingReceived) error
}
