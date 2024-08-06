package logger

import (
	"context"
	"fmt"
	"io"
	"log"
	"log/slog"
	"strings"
)

type ContextHandlerOptions struct {
	SlogOpts slog.HandlerOptions
}

type ContextHandler struct {
	slog.Handler
	l *log.Logger
}

// Handle processes and formats a log record into a structured log message
func (h *ContextHandler) Handle(ctx context.Context, r slog.Record) error {
	levelStr := r.Level.String()

	fields, err := jsonMarshal(formatAttributes(r))
	if err != nil {
		return err
	}

	timeStr := r.Time.Format("2006-01-02T15:04:05.000Z07:00") // ISO8601 format for timestamp
	msg := r.Message

	var builder strings.Builder
	builder.WriteString(fmt.Sprintf("%s [%s] %s %s", timeStr, levelStr, msg, string(fields)))

	h.l.Println(builder.String())
	return nil
}

// NewContextHandler creates a new ContextHandler with provided options
func NewContextHandler(out io.Writer, opts ContextHandlerOptions) slog.Handler {
	return &ContextHandler{
		Handler: slog.NewJSONHandler(out, &opts.SlogOpts),
		l:       log.New(out, "", 0),
	}
}
