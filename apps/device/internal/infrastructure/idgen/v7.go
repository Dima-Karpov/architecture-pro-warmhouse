package idgen

import "github.com/google/uuid"

type V7 struct{}

func NewV7() *V7 {
	return &V7{}
}

func (V7) New() uuid.UUID {
	id, err := uuid.NewV7()
	if err != nil {
		return uuid.New()
	}

	return id
}
