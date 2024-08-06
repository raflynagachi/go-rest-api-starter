package logger

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"log/slog"
	"testing"

	"github.com/fatih/color"
)

func TestPrettyHandler_Handle(t *testing.T) {

	t.Run("Success", func(t *testing.T) {
		var buf bytes.Buffer
		handler := NewPrettyHandler(&buf, PrettyHandlerOptions{})

		record := createTestRecord(slog.LevelInfo, "test message", slog.String("key", "value"))

		err := handler.Handle(context.Background(), record)
		if err != nil {
			t.Fatalf("Handle() returned an error: %v", err)
		}

		output := buf.String()

		timeStr := record.Time.Format("[15:05:05.000]")
		levelStr := formatLevel(record.Level)
		message := record.Message
		expectedFields, _ := json.MarshalIndent(formatAttributes(record), "", "  ")

		expected := fmt.Sprintf("%s %s %s %s\n", timeStr, levelStr, message, color.WhiteString(string(expectedFields)))

		if output != expected {
			t.Errorf("Handle() output = %q, want %q", output, expected)
		}
	})

	t.Run("MarshalError", func(t *testing.T) {
		var buf bytes.Buffer
		handler := &PrettyHandler{
			Handler: slog.NewJSONHandler(&buf, &slog.HandlerOptions{}),
			l:       log.New(&buf, "", 0),
		}

		originalMarshalIndent := json.MarshalIndent
		defer func() { jsonMarshalIndent = originalMarshalIndent }()
		jsonMarshalIndent = func(v any, prefix, indent string) ([]byte, error) {
			return nil, errors.New("mock marshaling error")
		}

		record := createTestRecord(slog.LevelInfo, "test message", slog.String("key", "value"))

		err := handler.Handle(context.Background(), record)
		if err == nil {
			t.Fatalf("Handle() did not return an error")
		}

		expectedErr := "mock marshaling error"
		if err.Error() != expectedErr {
			t.Errorf("Handle() error = %v, want %v", err, expectedErr)
		}
	})
}

func TestFormatLevel(t *testing.T) {
	tests := []struct {
		level    slog.Level
		expected string
	}{
		{slog.LevelDebug, color.MagentaString("DEBUG:")},
		{slog.LevelInfo, color.BlueString("INFO:")},
		{slog.LevelWarn, color.YellowString("WARN:")},
		{slog.LevelError, color.RedString("ERROR:")},
		{slog.Level(991), slog.Level(991).String() + ":"},
	}

	for _, tt := range tests {
		t.Run(tt.level.String(), func(t *testing.T) {
			result := formatLevel(tt.level)
			if result != tt.expected {
				t.Errorf("formatLevel(%v) = %q, want %q", tt.level, result, tt.expected)
			}
		})
	}
}
