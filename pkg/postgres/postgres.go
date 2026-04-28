package postgres

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/url"
	"strconv"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/prometheus/client_golang/prometheus"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"

	"github.com/vladimirkorzhenevskiy/gogen/pkg/log"
)

type Config struct {
	Host     string `envconfig:"HOST"      required:"false" default:"localhost"`
	Port     uint16 `envconfig:"PORT"      required:"false" default:"5432"`
	Username string `envconfig:"USERNAME"  required:"true"`
	Password string `envconfig:"PASSWORD"  required:"true"`
	DB       string `envconfig:"DB"        required:"true"`
	SSLMode  string `envconfig:"SSL_MODE"  required:"false" default:"disable"`
	MaxConns int32  `envconfig:"MAX_CONNS" required:"false" default:"10"`
}

func (c Config) DSN() string {
	return fmt.Sprintf("postgres://%s:%s@%s/%s?sslmode=%s",
		c.Username,
		url.QueryEscape(c.Password),
		net.JoinHostPort(c.Host, strconv.Itoa(int(c.Port))),
		c.DB,
		c.SSLMode,
	)
}

type Driver struct {
	*pgxpool.Pool

	metrics Metrics
	logger  *log.Logger
	tracer  trace.Tracer
}

func New(cfg Config) (*Driver, error) {
	poolConfig, err := pgxpool.ParseConfig(cfg.DSN())
	if err != nil {
		return nil, fmt.Errorf("postgres: failed to parse dsn: %w", err)
	}

	driver := &Driver{
		logger: log.Nop(),
		tracer: otel.GetTracerProvider().Tracer("github.com/vladimirkorzhenevskiy/gogen"),
	}

	poolConfig.MaxConns = cfg.MaxConns
	poolConfig.ConnConfig.Tracer = driver

	const connectTimeout = 10 * time.Second

	ctx, cancel := context.WithTimeout(context.Background(), connectTimeout)
	defer cancel()

	pool, err := pgxpool.NewWithConfig(ctx, poolConfig)
	if err != nil {
		return nil, fmt.Errorf("postgres: failed to open pgxpool connect: %w", err)
	}

	if err = pool.Ping(ctx); err != nil {
		return nil, fmt.Errorf("postgres: failed to ping pgxpool connect: %w", err)
	}

	driver.Pool = pool

	return driver, nil
}

func (d *Driver) WithLogger(logger *log.Logger) {
	d.logger = logger
}

func (d *Driver) WithTracer(tracer trace.Tracer) {
	d.tracer = tracer
}

func (d *Driver) TraceQueryStart(ctx context.Context, _ *pgx.Conn, _ pgx.TraceQueryStartData) context.Context {
	ctx, _ = d.tracer.Start(ctx, "postgres."+Op(ctx))

	return withQueryStart(ctx, time.Now())
}

func (d *Driver) TraceQueryEnd(ctx context.Context, _ *pgx.Conn, data pgx.TraceQueryEndData) {
	span := trace.SpanFromContext(ctx)
	defer span.End()

	op := Op(ctx)
	start := queryStart(ctx)
	duration := time.Since(start)

	if data.Err != nil && !errors.Is(data.Err, pgx.ErrNoRows) {
		d.logger.ErrorContext(ctx, "postgres: query failed",
			log.String("operation", op),
			log.Duration("duration", duration),
			log.Error(data.Err),
		)

		span.RecordError(data.Err)
	}

	labels := prometheus.Labels{
		"operation": op,
		"status":    getStatus(data.Err),
	}

	if d.metrics.QueriesTotal != nil {
		d.metrics.QueriesTotal.With(labels).Inc()
	}

	if d.metrics.QueriesDuration != nil {
		d.metrics.QueriesDuration.With(labels).Observe(duration.Seconds())
	}
}
