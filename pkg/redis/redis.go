package redis

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/redis/go-redis/v9"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"

	"github.com/vladimirkorzhenevskiy/gogen/pkg/log"
)

type Config struct {
	Host     string        `envconfig:"HOST"     required:"false" default:"localhost"`
	Port     uint16        `envconfig:"PORT"     required:"false" default:"6379"`
	Username string        `envconfig:"USERNAME" required:"false"`
	Password string        `envconfig:"PASSWORD" required:"false"`
	DB       int           `envconfig:"DB"       required:"false" default:"0"`
	Timeout  time.Duration `envconfig:"TIMEOUT"  required:"false" default:"1s"`
}

func (c Config) Address() string {
	return fmt.Sprintf("%s:%d", c.Host, c.Port)
}

type Driver struct {
	*redis.Client

	cfg     Config
	metrics Metrics
	logger  *log.Logger
	tracer  trace.Tracer
}

func New(cfg Config) (*Driver, error) {
	client := redis.NewClient(&redis.Options{
		Addr:         cfg.Address(),
		Username:     cfg.Username,
		Password:     cfg.Password,
		DB:           cfg.DB,
		ReadTimeout:  cfg.Timeout,
		WriteTimeout: cfg.Timeout,
	})

	const connectTimeout = 5 * time.Second

	ctx, cancel := context.WithTimeout(context.Background(), connectTimeout)
	defer cancel()

	if err := client.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("redis: failed to ping connect: %w", err)
	}

	driver := &Driver{
		Client: client,
		cfg:    cfg,
		logger: log.Nop(),
		tracer: otel.GetTracerProvider().Tracer("github.com/vladimirkorzhenevskiy/gogen"),
	}

	driver.AddHook(driver)

	return driver, nil
}

func (d *Driver) WithLogger(logger *log.Logger) {
	d.logger = logger
}

func (d *Driver) WithTracer(tracer trace.Tracer) {
	d.tracer = tracer
}

func (d *Driver) WithCommandsTotalCounter(counter *prometheus.CounterVec) {
	d.metrics.CommandsTotal = counter
}

func (d *Driver) WithCommandsDurationObserver(observer *prometheus.HistogramVec) {
	d.metrics.CommandsDuration = observer
}

func (d *Driver) WithPipelinesTotalCounter(counter *prometheus.CounterVec) {
	d.metrics.PipelinesTotal = counter
}

func (d *Driver) WithPipelinesCommandsTotalCounter(counter *prometheus.CounterVec) {
	d.metrics.PipelinesCommandsTotal = counter
}

func (d *Driver) WithPipelinesDurationObserver(observer *prometheus.HistogramVec) {
	d.metrics.PipelinesDuration = observer
}

func (d *Driver) DialHook(next redis.DialHook) redis.DialHook {
	return next
}

func (d *Driver) ProcessHook(next redis.ProcessHook) redis.ProcessHook {
	return func(ctx context.Context, command redis.Cmder) error {
		op := Op(ctx)

		ctx, span := d.tracer.Start(ctx, "redis.command."+op)
		defer span.End()

		start := time.Now()

		err := next(ctx, command)

		duration := time.Since(start)

		if err != nil && !errors.Is(err, redis.Nil) {
			d.logger.ErrorContext(ctx, "redis: command failed",
				log.String("operation", op),
				log.String("command", command.Name()),
				log.Duration("duration", duration),
				log.Error(err),
			)

			span.RecordError(err)
		}

		labels := prometheus.Labels{
			"operation": op,
			"command":   command.Name(),
			"status":    getStatus(err),
		}

		if d.metrics.CommandsTotal != nil {
			d.metrics.CommandsTotal.With(labels).Inc()
		}

		if d.metrics.CommandsDuration != nil {
			d.metrics.CommandsDuration.With(labels).Observe(duration.Seconds())
		}

		return err
	}
}

func (d *Driver) ProcessPipelineHook(next redis.ProcessPipelineHook) redis.ProcessPipelineHook {
	return func(ctx context.Context, commands []redis.Cmder) error {
		op := Op(ctx)

		ctx, span := d.tracer.Start(ctx, "redis.pipeline."+op)
		defer span.End()

		start := time.Now()

		err := next(ctx, commands)

		duration := time.Since(start)

		if err != nil && !errors.Is(err, redis.Nil) {
			d.logger.ErrorContext(ctx, "redis: pipeline failed",
				log.String("operation", op),
				log.Duration("duration", duration),
				log.Error(err),
			)

			span.RecordError(err)
		}

		labels := prometheus.Labels{
			"operation": op,
			"status":    getStatus(err),
		}

		if d.metrics.PipelinesTotal != nil {
			d.metrics.PipelinesTotal.With(labels).Inc()
		}

		if d.metrics.PipelinesCommandsTotal != nil {
			d.metrics.PipelinesCommandsTotal.With(labels).Inc()
		}

		if d.metrics.PipelinesDuration != nil {
			d.metrics.PipelinesDuration.With(labels).Observe(duration.Seconds())
		}

		return err
	}
}
