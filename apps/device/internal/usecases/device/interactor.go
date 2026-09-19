package device

import (
	"context"
	"errors"

	"github.com/google/uuid"

	"device/internal/domain"
)

type Interactor struct {
	repo    Repository
	ids     IDGenerator
	houseID uuid.UUID
	typeID  uuid.UUID
}

func NewInteractor(repo Repository, ids IDGenerator, houseID, typeID uuid.UUID) *Interactor {
	return &Interactor{
		repo:    repo,
		ids:     ids,
		houseID: houseID,
		typeID:  typeID,
	}
}

func (interactor *Interactor) Get(ctx context.Context, id uuid.UUID) (domain.Device, error) {
	return interactor.repo.Get(ctx, id)
}

func (interactor *Interactor) List(ctx context.Context) ([]domain.Device, error) {
	return interactor.repo.List(ctx)
}

func (interactor *Interactor) UpdateStatus(
	ctx context.Context,
	id uuid.UUID,
	status string,
) (domain.Device, error) {
	parsed, err := domain.ParseStatus(status)
	if err != nil {
		return domain.Device{}, err
	}

	item, err := interactor.repo.Get(ctx, id)
	if err != nil {
		return domain.Device{}, err
	}

	item.Status = parsed

	return interactor.repo.Save(ctx, item)
}

type UpsertInput struct {
	SerialNumber string
	Address      string
	TypeCode     string
	Status       string
	ExternalID   string
}

func (interactor *Interactor) Upsert(ctx context.Context, in UpsertInput) (domain.Device, error) {
	if in.SerialNumber == "" {
		return domain.Device{}, domain.ErrInvalidPayload
	}

	typeCode, err := domain.ParseTypeCode(in.TypeCode)
	if err != nil {
		return domain.Device{}, err
	}

	status, err := domain.ParseStatus(in.Status)
	if err != nil {
		return domain.Device{}, err
	}

	item := domain.Device{
		SerialNumber: in.SerialNumber,
		Address:      in.Address,
		TypeCode:     typeCode,
		Status:       status,
		ExternalID:   in.ExternalID,
		TypeID:       interactor.typeID,
		HouseID:      interactor.houseID,
	}

	if in.ExternalID != "" {
		existing, getErr := interactor.repo.GetByExternalID(ctx, in.ExternalID)
		if getErr == nil {
			item.ID = existing.ID
			item.TypeID = existing.TypeID
			item.HouseID = existing.HouseID

			return interactor.repo.Save(ctx, item)
		}

		if !errors.Is(getErr, domain.ErrNotFound) {
			return domain.Device{}, getErr
		}
	}

	item.ID = interactor.ids.New()

	return interactor.repo.Save(ctx, item)
}
