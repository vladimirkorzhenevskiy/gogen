package postgres

import (
	"context"
	"errors"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/prometheus/client_golang/prometheus"
)

type Metrics struct {
	QueriesTotal    *prometheus.CounterVec
	QueriesDuration *prometheus.HistogramVec
}

func getStatus(err error) string {
	switch {
	case err == nil || errors.Is(err, pgx.ErrNoRows):
		return "ok"
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

type queryStartKey struct{}

func queryStart(ctx context.Context) time.Time {
	v, ok := ctx.Value(queryStartKey{}).(time.Time)
	if !ok {
		return time.Now()
	}

	return v
}

func withQueryStart(ctx context.Context, startTime time.Time) context.Context {
	return context.WithValue(ctx, queryStartKey{}, startTime)
}
