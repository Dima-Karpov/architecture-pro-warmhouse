package temperature

import "temperature-api/internal/domain"

type Interactor struct {
	clock  Clock
	values ValueSource
}

func NewInteractor(clock Clock, values ValueSource) *Interactor {
	return &Interactor{
		clock:  clock,
		values: values,
	}
}

func (interactor *Interactor) Get(location, sensorID string) domain.Reading {
	return domain.NewReading(location, sensorID, interactor.values.Next(), interactor.clock.Now())
}
