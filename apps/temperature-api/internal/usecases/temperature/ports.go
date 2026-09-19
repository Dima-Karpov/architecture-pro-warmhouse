package temperature

import "time"

type Clock interface {
	Now() time.Time
}

type ValueSource interface {
	Next() float64
}
