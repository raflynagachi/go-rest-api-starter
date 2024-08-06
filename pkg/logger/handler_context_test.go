package logger

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"testing"
)

func TestContextHandler_Handle(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		var buf bytes.Buffer
		handler := NewContextHandler(&buf, ContextHandlerOptions{})

		record := createTestRecord(slog.LevelInfo, "test message", slog.String("key", "value"))

		err := handler.Handle(context.Background(), record)
		if err != nil {
			t.Fatalf("Handle() returned an error: %v", err)
		}

		output := buf.String()

		timeStr := record.Time.Format("2006-01-02T15:04:05.000Z07:00")
		levelStr := record.Level.String()
		message := record.Message
		expectedFields, _ := json.Marshal(formatAttributes(record))

		expected := fmt.Sprintf("%s [%s] %s %s\n", timeStr, levelStr, message, expectedFields)

		if output != expected {
			t.Errorf("Handle() output = %q, want %q", output, expected)
		}
	})

	t.Run("MarshalError", func(t *testing.T) {
		var buf bytes.Buffer
		handler := NewContextHandler(&buf, ContextHandlerOptions{})

		originalMarshal := jsonMarshal
		defer func() { jsonMarshal = originalMarshal }()
		jsonMarshal = func(v any) ([]byte, error) {
			return nil, errors.New("test jsonMarshal error")
		}

		record := createTestRecord(slog.LevelInfo, "test message", slog.String("key", "value"))

		err := handler.Handle(context.Background(), record)
		if err == nil {
			t.Fatalf("Handle() did not return an error")
		}

		expectedErr := "test jsonMarshal error"
		if err.Error() != expectedErr {
			t.Errorf("Handle() error = %v, want %v", err, expectedErr)
		}
	})
}
