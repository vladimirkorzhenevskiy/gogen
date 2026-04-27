package log

import (
	"context"
	"io"
	"log/slog"
	"os"
)

type Logger struct {
	logger *slog.Logger
}

func New(opts ...Option) *Logger {
	options := Options{
		Writer: os.Stderr,
		Level:  LevelFromEnv(),
	}

	for _, opt := range opts {
		opt(&options)
	}

	handler := newHandler(options.Writer, options.Level)

	return &Logger{logger: slog.New(handler)}
}

func Nop() *Logger {
	return &Logger{logger: slog.New(slog.DiscardHandler)}
}

func (l *Logger) Debug(msg string, args ...any) {
	l.logger.Debug(msg, args...)
}

func (l *Logger) DebugContext(ctx context.Context, msg string, args ...any) {
	l.logger.DebugContext(ctx, msg, append(ContextArgs(ctx), args...)...)
}

func (l *Logger) Info(msg string, args ...any) {
	l.logger.Info(msg, args...)
}

func (l *Logger) InfoContext(ctx context.Context, msg string, args ...any) {
	l.logger.InfoContext(ctx, msg, append(ContextArgs(ctx), args...)...)
}

func (l *Logger) Warn(msg string, args ...any) {
	l.logger.Warn(msg, args...)
}

func (l *Logger) WarnContext(ctx context.Context, msg string, args ...any) {
	l.logger.WarnContext(ctx, msg, append(ContextArgs(ctx), args...)...)
}

func (l *Logger) Error(msg string, args ...any) {
	l.logger.Error(msg, args...)
}

func (l *Logger) ErrorContext(ctx context.Context, msg string, args ...any) {
	l.logger.ErrorContext(ctx, msg, append(ContextArgs(ctx), args...)...)
}

func (l *Logger) Fatal(msg string, args ...any) {
	l.logger.Error(msg, args...)
	os.Exit(1)
}

func (l *Logger) FatalContext(ctx context.Context, msg string, args ...any) {
	l.logger.ErrorContext(ctx, msg, args...)
	os.Exit(1)
}

func (l *Logger) Handler() slog.Handler {
	return l.logger.Handler()
}

func newHandler(writer io.Writer, level Level) *Handler {
	return &Handler{
		Handler: slog.NewJSONHandler(writer, &slog.HandlerOptions{
			Level: level,
		}),
	}
}

type Handler struct {
	slog.Handler
}

func (h *Handler) Handle(ctx context.Context, record slog.Record) error {
	record.Add(ContextArgs(ctx)...)

	return h.Handle(ctx, record)
}
