package postgres

import (
	"context"
	"fmt"
	"net"
	"net/url"
	"strconv"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
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
}

func New(cfg Config) (*Driver, error) {
	poolConfig, err := pgxpool.ParseConfig(cfg.DSN())
	if err != nil {
		return nil, fmt.Errorf("db: failed to parse dsn: %w", err)
	}

	poolConfig.MaxConns = cfg.MaxConns

	const connectTimeout = 10 * time.Second

	ctx, cancel := context.WithTimeout(context.Background(), connectTimeout)
	defer cancel()

	pool, err := pgxpool.New(ctx, cfg.DSN())
	if err != nil {
		return nil, fmt.Errorf("db: failed to open connect: %w", err)
	}

	if err = pool.Ping(ctx); err != nil {
		return nil, fmt.Errorf("db: failed to ping connect: %w", err)
	}

	return &Driver{Pool: pool}, nil
}

func (m *Driver) WithLogger(logger *log.Logger) {

}

func (m *Driver) WithTracer(tracer trace.Tracer) {

}
