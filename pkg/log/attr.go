package log

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"time"

	"github.com/google/uuid"
)

type Attr = slog.Attr

func Any(key string, value any) Attr {
	str, _ := json.Marshal(value) //nolint:errchkjson

	return slog.String(key, string(str))
}

func String(key, value string) Attr {
	return slog.String(key, value)
}

func Stringer(key string, value fmt.Stringer) Attr {
	return slog.String(key, value.String())
}

func Bytes(key string, value []byte) Attr {
	return slog.String(key, string(value))
}

func Bool(key string, value bool) Attr {
	return slog.Bool(key, value)
}

func Int(key string, value int) Attr {
	return slog.Int(key, value)
}

func Int64(key string, value int64) Attr {
	return slog.Int64(key, value)
}

func Uint64(key string, value uint64) Attr {
	return slog.Uint64(key, value)
}

func Float64(key string, value float64) Attr {
	return slog.Float64(key, value)
}

func UUID(key string, value uuid.UUID) Attr {
	return slog.String(key, value.String())
}

func Time(key string, value time.Time) Attr {
	return slog.Time(key, value)
}

func Duration(key string, value time.Duration) Attr {
	return slog.Duration(key, value)
}

func Error(err error) Attr {
	return slog.String("err", err.Error())
}

func Group(key string, values ...any) Attr {
	return slog.Group(key, values...)
}
