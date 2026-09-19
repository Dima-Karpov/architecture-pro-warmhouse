package responses

import (
	"time"

	"temperature-api/internal/domain"
)

type TemperatureResponse struct {
	Timestamp   time.Time `json:"timestamp"`
	Unit        string    `json:"unit"`
	Location    string    `json:"location"`
	Status      string    `json:"status"`
	SensorID    string    `json:"sensor_id"`
	SensorType  string    `json:"sensor_type"`
	Description string    `json:"description"`
	Value       float64   `json:"value"`
}

func NewTemperatureResponse(reading domain.Reading) TemperatureResponse {
	return TemperatureResponse{
		Timestamp:   reading.Timestamp,
		Unit:        reading.Unit,
		Location:    reading.Location,
		Status:      reading.Status,
		SensorID:    reading.SensorID,
		SensorType:  reading.SensorType,
		Description: reading.Description,
		Value:       reading.Value,
	}
}
