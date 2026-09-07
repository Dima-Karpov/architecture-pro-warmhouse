package api

import "temperature-api/internal/domain"

type TemperatureGetter interface {
	Get(location, sensorID string) domain.Reading
}
