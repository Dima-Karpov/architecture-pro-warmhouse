package device

import (
	"context"
	"testing"

	"github.com/google/uuid"

	"device/internal/domain"
)

type stubRepo struct {
	byID       map[uuid.UUID]domain.Device
	byExternal map[string]domain.Device
}

func newStubRepo() *stubRepo {
	return &stubRepo{
		byID:       make(map[uuid.UUID]domain.Device),
		byExternal: make(map[string]domain.Device),
	}
}

func (s *stubRepo) Get(_ context.Context, id uuid.UUID) (domain.Device, error) {
	item, ok := s.byID[id]
	if !ok {
		return domain.Device{}, domain.ErrNotFound
	}

	return item, nil
}

func (s *stubRepo) GetByExternalID(_ context.Context, externalID string) (domain.Device, error) {
	item, ok := s.byExternal[externalID]
	if !ok {
		return domain.Device{}, domain.ErrNotFound
	}

	return item, nil
}

func (s *stubRepo) List(_ context.Context) ([]domain.Device, error) {
	items := make([]domain.Device, 0, len(s.byID))
	for _, item := range s.byID {
		items = append(items, item)
	}

	return items, nil
}

func (s *stubRepo) Save(_ context.Context, item domain.Device) (domain.Device, error) {
	s.byID[item.ID] = item
	if item.ExternalID != "" {
		s.byExternal[item.ExternalID] = item
	}

	return item, nil
}

type stubIDs struct {
	id uuid.UUID
}

func (s stubIDs) New() uuid.UUID {
	return s.id
}

func TestInteractorUpsertCreatesThenUpdates(t *testing.T) {
	t.Parallel()

	houseID := uuid.MustParse("01932c4e-8a1d-7222-8333-3c6d7e8f9012")
	typeID := uuid.MustParse("01932c4e-8a1c-7111-8222-2b5c6d7e8f90")
	newID := uuid.MustParse("01932c4e-8a1b-7f3c-8d2e-1a4b5c6d7e8f")
	interactor := NewInteractor(newStubRepo(), stubIDs{id: newID}, houseID, typeID)
	ctx := context.Background()

	created, err := interactor.Upsert(ctx, UpsertInput{
		SerialNumber: "RL-22",
		Address:      "Living Room",
		TypeCode:     domain.TypeHeating,
		Status:       domain.StatusOn,
		ExternalID:   "1",
	})
	if err != nil {
		t.Fatalf("upsert create: %v", err)
	}

	if created.ID != newID {
		t.Fatalf("id: got %s, want %s", created.ID, newID)
	}

	updated, err := interactor.Upsert(ctx, UpsertInput{
		SerialNumber: "RL-22",
		Address:      "Bedroom",
		TypeCode:     domain.TypeHeating,
		Status:       domain.StatusOff,
		ExternalID:   "1",
	})
	if err != nil {
		t.Fatalf("upsert update: %v", err)
	}

	if updated.ID != newID {
		t.Fatalf("updated id changed: %s", updated.ID)
	}

	if updated.Address != "Bedroom" || updated.Status != domain.StatusOff {
		t.Fatalf("update not applied: %+v", updated)
	}
}

func TestInteractorUpdateStatusNotFound(t *testing.T) {
	t.Parallel()

	interactor := NewInteractor(newStubRepo(), stubIDs{}, uuid.Nil, uuid.Nil)
	_, err := interactor.UpdateStatus(
		context.Background(),
		uuid.MustParse("01932c4e-8a1b-7f3c-8d2e-1a4b5c6d7e8f"),
		domain.StatusOff,
	)

	if err == nil {
		t.Fatal("expected not found")
	}
}
