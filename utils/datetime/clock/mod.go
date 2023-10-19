package clock

import "time"

func Fixed() *Clock {
	return nil
}

func Offset(base *Clock, offset time.Duration) *Clock {
	return nil
}

type Clock struct {
}

func (c *Clock) millis() int64 {
	return 0
}
