package device

import (
	"context"

	"github.com/google/uuid"

	"device/internal/domain"
)

type Repository interface {
	Get(ctx context.Context, id uuid.UUID) (domain.Device, error)
	GetByExternalID(ctx context.Context, externalID string) (domain.Device, error)
	List(ctx context.Context) ([]domain.Device, error)
	Save(ctx context.Context, item domain.Device) (domain.Device, error)
}

type IDGenerator interface {
	New() uuid.UUID
}
