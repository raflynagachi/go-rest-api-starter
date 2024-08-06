package logger

import (
	"testing"
)

func TestWithLevel(t *testing.T) {
	tests := []struct {
		name      string
		level     Level
		wantLevel Level
	}{
		{name: "Debug Level", level: LevelDebug, wantLevel: LevelDebug},
		{name: "Info Level", level: LevelInfo, wantLevel: LevelInfo},
		{name: "Error Level", level: LevelError, wantLevel: LevelError},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			opts := &Options{}
			WithLevel(tt.level)(opts)
			if opts.Level != tt.wantLevel {
				t.Errorf("WithLevel() = %v, want %v", opts.Level, tt.wantLevel)
			}
		})
	}
}

func TestWithAddSource(t *testing.T) {
	tests := []struct {
		name          string
		addSource     bool
		wantAddSource bool
	}{
		{name: "Add Source True", addSource: true, wantAddSource: true},
		{name: "Add Source False", addSource: false, wantAddSource: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			opts := &Options{}
			WithAddSource(tt.addSource)(opts)
			if opts.AddSource != tt.wantAddSource {
				t.Errorf("WithAddSource() = %v, want %v", opts.AddSource, tt.wantAddSource)
			}
		})
	}
}

func TestWithEnv(t *testing.T) {
	tests := []struct {
		name    string
		env     string
		wantEnv string
	}{
		{name: "Development Env", env: "development", wantEnv: "development"},
		{name: "Production Env", env: "production", wantEnv: "production"},
		{name: "Empty Env", env: "", wantEnv: ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			opts := &Options{}
			WithEnv(tt.env)(opts)
			if opts.Env != tt.wantEnv {
				t.Errorf("WithEnv() = %v, want %v", opts.Env, tt.wantEnv)
			}
		})
	}
}
