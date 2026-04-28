package redis

import (
	"context"
	"errors"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/redis/go-redis/v9"
)

type Metrics struct {
	CommandsTotal          *prometheus.CounterVec
	CommandsDuration       *prometheus.HistogramVec
	PipelinesTotal         *prometheus.CounterVec
	PipelinesCommandsTotal *prometheus.CounterVec
	PipelinesDuration      *prometheus.HistogramVec
}

func getStatus(err error) string {
	switch {
	case err == nil:
		return "ok"
	case errors.Is(err, redis.Nil):
		return "miss"
	default:
		return "error"
	}
}

type operationKey struct{}

func Op(ctx context.Context) string {
	if v, ok := ctx.Value(operationKey{}).(string); ok {
		return v
	}

	return "unknown"
}

func WithOp(ctx context.Context, name string) context.Context {
	return context.WithValue(ctx, operationKey{}, name)
}
