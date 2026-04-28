package server

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net"
	"net/http"
	"sync"
	"time"

	"github.com/vladimirkorzhenevskiy/gogen/pkg/log"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
)

var ErrAlreadyRunning = errors.New("server is already running")

type Config struct {
	Address      string        `envconfig:"ADDRESS"       default:"0.0.0.0:8080"`
	ReadTimeout  time.Duration `envconfig:"READ_TIMEOUT"  default:"5s"`
	WriteTimeout time.Duration `envconfig:"WRITE_TIMEOUT" default:"60s"`
}

type Server struct {
	cfg         Config
	server      *http.Server
	handler     http.Handler
	middlewares []Middleware
	logger      *log.Logger
	tracer      trace.Tracer
	mx          sync.Mutex
	running     bool
}

func New(cfg Config) (*Server, error) {
	return &Server{
		cfg:     cfg,
		handler: http.DefaultServeMux,
		logger:  log.Nop(),
		tracer:  otel.GetTracerProvider().Tracer("github.com/vladimirkorzhenevskiy/gogen"),
	}, nil
}

func (s *Server) WithLogger(logger *log.Logger) *Server {
	s.logger = logger

	return s
}

func (s *Server) WithTracer(tracer trace.Tracer) *Server {
	s.tracer = tracer

	return s
}

func (s *Server) WithHandler(handler http.Handler) *Server {
	if handler == nil {
		s.handler = http.DefaultServeMux
	} else {
		s.handler = handler
	}

	return s
}

func (s *Server) WithMiddlewares(middlewares ...Middleware) *Server {
	s.middlewares = append(s.middlewares, middlewares...)

	return s
}

func (s *Server) Start(ctx context.Context) error {
	s.mx.Lock()
	defer s.mx.Unlock()

	if s.running {
		return ErrAlreadyRunning
	}

	listener, err := s.listen(ctx, s.cfg.Address)
	if err != nil {
		return err
	}

	s.setup()
	s.running = true
	s.run(listener)

	return nil
}

func (s *Server) Stop(ctx context.Context) error {
	s.mx.Lock()
	defer s.mx.Unlock()

	if !s.running {
		return nil
	}

	s.running = false

	return s.server.Shutdown(ctx)
}

func (s *Server) setup() {
	s.server = &http.Server{
		Addr:         s.cfg.Address,
		ReadTimeout:  s.cfg.ReadTimeout,
		WriteTimeout: s.cfg.WriteTimeout,
		Handler:      s.handler,
		ErrorLog:     slog.NewLogLogger(s.logger.Handler(), slog.LevelError),
	}
}

func (s *Server) listen(ctx context.Context, address string) (net.Listener, error) {
	var cfg net.ListenConfig

	listener, err := cfg.Listen(ctx, "tcp", address)
	if err != nil {
		return nil, fmt.Errorf("listen TCP: %w", err)
	}

	return listener, nil
}

func (s *Server) run(listener net.Listener) {
	s.logger.Info("server: starting",
		log.String("addr", s.server.Addr),
	)

	go func() {
		s.mx.Lock()
		running := s.running
		s.mx.Unlock()

		if !running {
			return
		}

		s.logger.Info("server: running",
			log.String("addr", listener.Addr().String()),
		)

		if err := s.server.Serve(listener); errors.Is(err, http.ErrServerClosed) {
			s.logger.Info("server: stopped",
				log.String("addr", s.server.Addr),
			)
		} else if err != nil {
			s.logger.Error("server: stopped with error",
				log.String("addr", s.server.Addr),
				log.Error(err),
			)
		}
	}()
}
