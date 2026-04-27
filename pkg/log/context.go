package log

import (
	"context"
)

type contextKey struct{}

func Context(ctx context.Context, args ...any) context.Context {
	if len(args) == 0 {
		return ctx
	}

	a, ok := ctx.Value(contextKey{}).([]any)
	if ok {
		a = append(a, args...)
	} else {
		a = args
	}

	return context.WithValue(ctx, contextKey{}, a)
}

func ContextArgs(ctx context.Context) []any {
	args, ok := ctx.Value(contextKey{}).([]any)
	if ok {
		return args
	}

	return nil
}
