package scenario

import "scenario/internal/domain"

type Interactor struct{}

func NewInteractor() *Interactor {
	return &Interactor{}
}

func (Interactor) Handle(event domain.ReadingReceived) domain.Decision {
	return domain.Decide(event.Value)
}
