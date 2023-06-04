package model

import "time"

type ITimeController interface {
	Now() time.Time
}

type TimeControllerFunc func() time.Time

func (t TimeControllerFunc) Now() time.Time {
	return t()
}
