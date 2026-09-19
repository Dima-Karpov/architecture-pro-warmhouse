package db

import (
	"time"

	"github.com/google/uuid"

	"device/internal/domain"
)

//nolint:govet // GORM-модель: много строк, fieldalignment не сходится
type deviceRow struct {
	ID           uuid.UUID `gorm:"type:uuid;primaryKey"`
	TypeID       uuid.UUID `gorm:"type:uuid;not null"`
	HouseID      uuid.UUID `gorm:"type:uuid;not null"`
	CreatedAt    time.Time `gorm:"autoCreateTime"`
	SerialNumber string    `gorm:"not null"`
	Address      string    `gorm:"not null;default:''"`
	TypeCode     string    `gorm:"not null"`
	Status       string    `gorm:"not null"`
	ExternalID   *string   `gorm:"uniqueIndex"`
}

func (deviceRow) TableName() string {
	return "devices"
}

func toRow(item domain.Device) deviceRow {
	row := deviceRow{
		SerialNumber: item.SerialNumber,
		Address:      item.Address,
		TypeCode:     item.TypeCode,
		Status:       item.Status,
		ID:           item.ID,
		TypeID:       item.TypeID,
		HouseID:      item.HouseID,
	}

	if item.ExternalID != "" {
		external := item.ExternalID
		row.ExternalID = &external
	}

	return row
}

func toDevice(row deviceRow) domain.Device {
	external := ""
	if row.ExternalID != nil {
		external = *row.ExternalID
	}

	return domain.Device{
		SerialNumber: row.SerialNumber,
		Address:      row.Address,
		TypeCode:     row.TypeCode,
		Status:       row.Status,
		ExternalID:   external,
		ID:           row.ID,
		TypeID:       row.TypeID,
		HouseID:      row.HouseID,
	}
}
