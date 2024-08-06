package logger

import (
	"io"
	"log/slog"
	"os"
)

func NewLogger(opts ...Option) *slog.Logger {
	config := &Options{
		Level:     slog.LevelInfo,
		AddSource: true,
	}

	for _, opt := range opts {
		opt(config)
	}

	options := &slog.HandlerOptions{
		AddSource: config.AddSource,
		Level:     config.Level,
	}

	h := NewPrettyHandler(os.Stdout, PrettyHandlerOptions{SlogOpts: *options})
	switch config.Env {
	case "production":
		h = NewContextHandler(os.Stdout, ContextHandlerOptions{SlogOpts: *options})
	case "test":
		h = NewContextHandler(io.Discard, ContextHandlerOptions{SlogOpts: *options})
	}

	logger := slog.New(h)
	slog.SetDefault(logger)

	return logger
}
