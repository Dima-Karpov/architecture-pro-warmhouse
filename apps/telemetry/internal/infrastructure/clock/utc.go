package clock

import "time"

type UTC struct{}

func NewUTC() *UTC {
	return &UTC{}
}

func (UTC) Now() time.Time {
	return time.Now().UTC()
}
