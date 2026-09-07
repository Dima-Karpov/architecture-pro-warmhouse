package random

import "math/rand/v2"

const (
	tempMinTenths  = 150
	tempSpanTenths = 201
	tenths         = 10
)

type Source struct{}

func NewSource() *Source {
	return &Source{}
}

func (Source) Next() float64 {
	return float64(tempMinTenths+rand.IntN(tempSpanTenths)) / tenths
}
