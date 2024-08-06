package logger

import (
	"encoding/json"
	"log/slog"
)

const (
	LevelInfo  = slog.LevelInfo
	LevelWarn  = slog.LevelWarn
	LevelError = slog.LevelError
	LevelDebug = slog.LevelDebug
)

type (
	Logger = slog.Logger
	Attr   = slog.Attr
	Level  = slog.Level
)

var (
	jsonMarshal       = json.Marshal
	jsonMarshalIndent = json.MarshalIndent
)
