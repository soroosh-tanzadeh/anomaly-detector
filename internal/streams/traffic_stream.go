package streams

import (
	"context"
	"time"
)

type TrafficStream interface {
	Add(context.Context, time.Time, int64) error
	Range(from time.Time, to time.Timer) []int64
}
