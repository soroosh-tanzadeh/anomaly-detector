package streams

import (
	"context"
	"time"
)

type TrafficStream interface {
	Add(context.Context, float64) error
	Range(ctx context.Context, from time.Time, to time.Time) ([]float64, error)
}
