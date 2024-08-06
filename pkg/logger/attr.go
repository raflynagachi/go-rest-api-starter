package logger

import (
	"log/slog"
	"time"
)

// formatAttributes converts slog attributes to map string of interface{}
func formatAttributes(r slog.Record) map[string]interface{} {
	fields := make(map[string]interface{}, r.NumAttrs())
	r.Attrs(func(a slog.Attr) bool {
		fields[a.Key] = a.Value.Any()
		return true
	})

	return fields
}

// GroupAttr creates a group for attributes
func GroupAttr(key string, attr ...any) Attr {
	return slog.Group(key, attr...)
}

// StringAttr creates an attribute from a string value.
func StringAttr(key string, val string) Attr {
	return slog.String(key, val)
}

// Float64Attr creates an attribute from a float64 value.
func Float64Attr(key string, val float64) Attr {
	return slog.Float64(key, val)
}

// Int64Attr creates an attribute from an int64 value.
func Int64Attr(key string, val int64) Attr {
	return slog.Int64(key, val)
}

// Uint64Attr creates an attribute from a uint64 value.
func Uint64Attr(key string, val uint64) Attr {
	return slog.Uint64(key, val)
}

// TimeAttr creates an attribute from a time.Time value.
// The format is RFC3339.
func TimeAttr(key string, t time.Time) Attr {
	return slog.String(key, t.Format(time.RFC3339))
}

// ErrAttr creates an attribute from an error.
func ErrAttr(err error) Attr {
	if err == nil {
		return slog.String("error", "no error")
	}
	return slog.String("error", err.Error())
}
