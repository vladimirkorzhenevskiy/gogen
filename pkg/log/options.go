package log //nolint:revive

import (
	"io"
)

type Options struct {
	Writer io.Writer
	Level  Level
}

type Option func(opt *Options)

func WithWriter(writer io.Writer) Option {
	return func(opt *Options) {
		opt.Writer = writer
	}
}

func WithLevel(level Level) Option {
	return func(opt *Options) {
		opt.Level = level
	}
}
