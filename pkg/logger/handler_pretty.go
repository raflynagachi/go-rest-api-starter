package logger

import (
	"context"
	"fmt"
	"io"
	"log"
	"log/slog"
	"strings"

	"github.com/fatih/color"
)

type PrettyHandlerOptions struct {
	SlogOpts slog.HandlerOptions
}

type PrettyHandler struct {
	slog.Handler
	l *log.Logger
}

func (h *PrettyHandler) Handle(ctx context.Context, r slog.Record) error {
	levelStr := formatLevel(r.Level)

	b, err := jsonMarshalIndent(formatAttributes(r), "", "  ")
	if err != nil {
		return err
	}
	fields := color.WhiteString(string(b))

	timeStr := r.Time.Format("[15:05:05.000]")
	msg := color.CyanString(r.Message)

	var builder strings.Builder
	builder.WriteString(fmt.Sprintf("%s %s %s %s", timeStr, levelStr, msg, fields))

	h.l.Println(builder.String())
	return nil
}

func formatLevel(level slog.Level) string {
	levelStr := level.String() + ":"
	switch level {
	case slog.LevelDebug:
		return color.MagentaString(levelStr)
	case slog.LevelInfo:
		return color.BlueString(levelStr)
	case slog.LevelWarn:
		return color.YellowString(levelStr)
	case slog.LevelError:
		return color.RedString(levelStr)
	default:
		return levelStr
	}
}

func NewPrettyHandler(out io.Writer, opts PrettyHandlerOptions) slog.Handler {
	h := &PrettyHandler{
		Handler: slog.NewJSONHandler(out, &opts.SlogOpts),
		l:       log.New(out, "", 0),
	}

	return h
}
