package domain

func LocationByID(sensorID string) string {
	switch sensorID {
	case IDLivingRoom:
		return LocationLivingRoom
	case IDBedroom:
		return LocationBedroom
	case IDKitchen:
		return LocationKitchen
	default:
		return LocationUnknown
	}
}

func IDByLocation(location string) string {
	switch location {
	case LocationLivingRoom:
		return IDLivingRoom
	case LocationBedroom:
		return IDBedroom
	case LocationKitchen:
		return IDKitchen
	default:
		return IDUnknown
	}
}

func Resolve(location, sensorID string) (string, string) {
	if location == "" {
		location = LocationByID(sensorID)
	}

	if sensorID == "" {
		sensorID = IDByLocation(location)
	}

	return location, sensorID
}
