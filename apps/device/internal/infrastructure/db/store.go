package db

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"device/internal/domain"
)

type Store struct {
	conn *gorm.DB
}

func NewStore(conn *gorm.DB) *Store {
	return &Store{conn: conn}
}

func (store *Store) Get(ctx context.Context, id uuid.UUID) (domain.Device, error) {
	var row deviceRow

	err := store.conn.
		WithContext(ctx).
		First(&row, "id = ?", id).
		Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return domain.Device{}, domain.ErrNotFound
	}

	if err != nil {
		return domain.Device{}, fmt.Errorf("get device: %w", err)
	}

	return toDevice(row), nil
}

func (store *Store) GetByExternalID(ctx context.Context, externalID string) (domain.Device, error) {
	var row deviceRow

	err := store.conn.
		WithContext(ctx).
		First(&row, "external_id = ?", externalID).
		Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return domain.Device{}, domain.ErrNotFound
	}

	if err != nil {
		return domain.Device{}, fmt.Errorf("get device: %w", err)
	}

	return toDevice(row), nil
}

func (store *Store) List(ctx context.Context) ([]domain.Device, error) {
	var rows []deviceRow

	err := store.conn.WithContext(ctx).Order("created_at").Find(&rows).Error
	if err != nil {
		return nil, fmt.Errorf("list devices: %w", err)
	}

	items := make([]domain.Device, 0, len(rows))
	for _, row := range rows {
		items = append(items, toDevice(row))
	}

	return items, nil
}

func (store *Store) Save(ctx context.Context, item domain.Device) (domain.Device, error) {
	row := toRow(item)

	if err := store.conn.
		WithContext(ctx).
		Save(&row).
		Error; err != nil {
		return domain.Device{}, fmt.Errorf("save device: %w", err)
	}

	return toDevice(row), nil
}
