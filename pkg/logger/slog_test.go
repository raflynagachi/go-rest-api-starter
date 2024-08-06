package logger

import (
	"context"
	"fmt"
	"log/slog"
	"testing"
)

func TestNewLogger(t *testing.T) {
	tests := []struct {
		name                string
		options             []Option
		expectedLevel       slog.Level
		expectedEnv         string
		expectedHandlerType string
	}{
		{
			name:                "Development Environment with PrettyHandler",
			options:             []Option{WithLevel(LevelDebug), WithEnv("development")},
			expectedLevel:       LevelDebug,
			expectedEnv:         "development",
			expectedHandlerType: "*logger.PrettyHandler",
		},
		{
			name:                "Production Environment with ContextHandler",
			options:             []Option{WithLevel(LevelInfo), WithEnv("production")},
			expectedLevel:       LevelInfo,
			expectedEnv:         "production",
			expectedHandlerType: "*logger.ContextHandler",
		},
		{
			name:                "Test Environment with ContextHandler",
			options:             []Option{WithLevel(LevelInfo), WithEnv("test")},
			expectedLevel:       LevelInfo,
			expectedEnv:         "test",
			expectedHandlerType: "*logger.ContextHandler",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			logger := NewLogger(tt.options...)

			if logger == nil {
				t.Fatal("Expected logger to be non-nil")
			}

			handler := logger.Handler()
			if handler == nil {
				t.Fatal("Expected handler to be non-nil")
			}

			handlerType := fmt.Sprintf("%T", handler)
			if handlerType != tt.expectedHandlerType {
				t.Errorf("Expected handler type %s, got %s", tt.expectedHandlerType, handlerType)
			}

			if !logger.Handler().Enabled(context.Background(), tt.expectedLevel) {
				t.Errorf("Expected log level %v, but its not enabled for the logger", tt.expectedLevel)
			}
		})
	}
}
