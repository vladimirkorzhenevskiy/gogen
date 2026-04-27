package redis

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
	"go.opentelemetry.io/otel/trace"

	"github.com/vladimirkorzhenevskiy/gogen/pkg/log"
)

type Config struct {
	Host     string `envconfig:"HOST"     required:"false" default:"localhost"`
	Port     uint16 `envconfig:"PORT"     required:"false" default:"6379"`
	Username string `envconfig:"USERNAME" required:"false"`
	Password string `envconfig:"PASSWORD" required:"false"`
	DB       int    `envconfig:"DB"       required:"false" default:"0"`
}

func (c Config) Address() string {
	return fmt.Sprintf("%s:%d", c.Host, c.Port)
}

type Driver struct {
	*redis.Client
}

func New(cfg Config) (*Driver, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     cfg.Address(),
		Username: cfg.Username,
		Password: cfg.Password,
		DB:       cfg.DB,
	})

	const connectTimeout = 5 * time.Second

	ctx, cancel := context.WithTimeout(context.Background(), connectTimeout)
	defer cancel()

	if err := client.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("redis: failed to ping connect: %w", err)
	}

	return &Driver{Client: client}, nil
}

func (m *Driver) WithLogger(logger *log.Logger) {

}

func (m *Driver) WithTracer(tracer trace.Tracer) {

}
