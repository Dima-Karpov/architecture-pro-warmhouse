package domain

import "time"

const (
	SensorTypeTemperature = "temperature"
	StatusActive          = "active"
	UnitCelsius           = "°C"

	LocationLivingRoom = "Living Room"
	LocationBedroom    = "Bedroom"
	LocationKitchen    = "Kitchen"
	LocationUnknown    = "Unknown"

	IDLivingRoom = "1"
	IDBedroom    = "2"
	IDKitchen    = "3"
	IDUnknown    = "0"
)

type Reading struct {
	Timestamp   time.Time
	Unit        string
	Location    string
	Status      string
	SensorID    string
	SensorType  string
	Description string
	Value       float64
}

func NewReading(location, sensorID string, value float64, at time.Time) Reading {
	location, sensorID = Resolve(location, sensorID)

	return Reading{
		Timestamp:   at,
		Unit:        UnitCelsius,
		Location:    location,
		Status:      StatusActive,
		SensorID:    sensorID,
		SensorType:  SensorTypeTemperature,
		Description: location + " temperature",
		Value:       value,
	}
}
