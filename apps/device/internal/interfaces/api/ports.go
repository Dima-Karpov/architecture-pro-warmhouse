package api

import (
	"context"

	"github.com/google/uuid"

	"device/internal/domain"
	usecase "device/internal/usecases/device"
)

type DeviceUsecase interface {
	Get(ctx context.Context, id uuid.UUID) (domain.Device, error)
	List(ctx context.Context) ([]domain.Device, error)
	UpdateStatus(ctx context.Context, id uuid.UUID, status string) (domain.Device, error)
	Upsert(ctx context.Context, in usecase.UpsertInput) (domain.Device, error)
}
