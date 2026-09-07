package api

import (
	"context"

	"github.com/google/uuid"

	"telemetry/internal/domain"
	usecase "telemetry/internal/usecases/telemetry"
)

type TelemetryUsecase interface {
	Ingest(ctx context.Context, in usecase.IngestInput) (domain.Reading, error)
	List(ctx context.Context, ids []uuid.UUID) ([]domain.Reading, error)
}
