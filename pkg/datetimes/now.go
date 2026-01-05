package datetimes

import "time"

type Now interface {
	Now() time.Time
}

type UTCNow struct{}

func NewUTCNow() *UTCNow {
	return &UTCNow{}
}

func (n *UTCNow) Now() time.Time {
	return time.Now().UTC()
}
