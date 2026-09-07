package db

import (
	"time"

	"github.com/google/uuid"

	"telemetry/internal/domain"
)

type readingRow struct {
	RecordedAt time.Time
	Unit       string    `gorm:"not null"`
	ID         uuid.UUID `gorm:"type:uuid;primaryKey"`
	DeviceID   uuid.UUID `gorm:"type:uuid;index;not null"`
	Value      float64   `gorm:"not null"`
}

func (readingRow) TableName() string {
	return "readings"
}

func toRow(item domain.Reading) readingRow {
	return readingRow{
		RecordedAt: item.RecordedAt,
		Unit:       item.Unit,
		ID:         item.ID,
		DeviceID:   item.DeviceID,
		Value:      item.Value,
	}
}

func toReading(row readingRow) domain.Reading {
	return domain.Reading{
		RecordedAt: row.RecordedAt,
		Unit:       row.Unit,
		ID:         row.ID,
		DeviceID:   row.DeviceID,
		Value:      row.Value,
	}
}
