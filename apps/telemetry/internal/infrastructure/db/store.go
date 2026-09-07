package db

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"telemetry/internal/domain"
)

type Store struct {
	conn *gorm.DB
}

func NewStore(conn *gorm.DB) *Store {
	return &Store{conn: conn}
}

func (store *Store) Save(ctx context.Context, reading domain.Reading) (domain.Reading, error) {
	row := toRow(reading)

	if err := store.conn.WithContext(ctx).Create(&row).Error; err != nil {
		return domain.Reading{}, fmt.Errorf("save reading: %w", err)
	}

	return toReading(row), nil
}

func (store *Store) ListByDeviceIDs(ctx context.Context, ids []uuid.UUID) ([]domain.Reading, error) {
	var rows []readingRow

	err := store.conn.WithContext(ctx).
		Where("device_id IN ?", ids).
		Order("recorded_at DESC").
		Find(&rows).Error
	if err != nil {
		return nil, fmt.Errorf("list readings: %w", err)
	}

	items := make([]domain.Reading, 0, len(rows))
	for _, row := range rows {
		items = append(items, toReading(row))
	}

	return items, nil
}
