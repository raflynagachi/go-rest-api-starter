package config

import (
	"errors"
	"os"
	"reflect"
	"testing"

	"github.com/raflynagachi/go-rest-api-starter/pkg/fs"
)

func TestReadJsonConfig(t *testing.T) {
	tests := []struct {
		name        string
		fileContent string
		setupFile   func() string
		want        *Config
		wantErr     bool
	}{
		{
			name: "success with valid config",
			setupFile: func() string {
				file, err := fs.CreateTempFile(`{"app": {"name":"my-service", "port":1234} }`, "config-*.json")
				if err != nil {
					t.Fatalf("Failed CreateTempFile(): %v", err)
				}
				return file.Name()
			},
			want: &Config{App: App{
				Name: "my-service",
				Port: 1234,
			}},
			wantErr: false,
		},
		{
			name: "success with missing field",
			setupFile: func() string {
				file, err := fs.CreateTempFile(`{"app": {"name":"my-service"} }`, "config-*.json")
				if err != nil {
					t.Fatalf("Failed CreateTempFile(): %v", err)
				}
				return file.Name()
			},
			want:    &Config{App: App{Name: "my-service"}},
			wantErr: false,
		},
		{
			name: "error due to invalid JSON",
			setupFile: func() string {
				file, err := fs.CreateTempFile(`{"app": {"name":"my-service", "port":1234 `, "config-*.json")
				if err != nil {
					t.Fatalf("Failed CreateTempFile(): %v", err)
				}
				return file.Name()
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "error due to empty file",
			setupFile: func() string {
				file, err := fs.CreateTempFile(``, "config-*.json")
				if err != nil {
					t.Fatalf("Failed CreateTempFile(): %v", err)
				}
				return file.Name()
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "error due to file does not exist",
			setupFile: func() string {
				return "/path/to/nonexistent/file.json"
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "error due to permission denied",
			setupFile: func() string {
				file, err := fs.CreateTempFile(``, "config-*.json")
				if err != nil {
					t.Fatalf("Failed CreateTempFile(): %v", err)
				}

				err = os.Chmod(file.Name(), 0000) // No permissions
				if err != nil {
					t.Fatalf("Failed to change file permissions: %v", err)
				}

				return file.Name()
			},
			want:    nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			filePath := tt.setupFile()
			defer os.Remove(filePath)

			got, err := ReadJsonConfig(filePath)
			if (err != nil) != tt.wantErr {
				t.Errorf("ReadJsonConfig() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ReadJsonConfig() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestLoadConfig(t *testing.T) {
	tmpReadJsonConfig := readJsonConfig
	defer func() {
		readJsonConfig = tmpReadJsonConfig
	}()

	tests := []struct {
		name    string
		env     string
		setup   func()
		want    *Config
		wantErr bool
	}{
		{
			name: "success with production env",
			env:  Production,
			setup: func() {
				readJsonConfig = func(path string) (*Config, error) {
					return &Config{App: App{Name: "my-service"}}, nil
				}
			},
			want:    &Config{App: App{Name: "my-service"}},
			wantErr: false,
		},
		{
			name: "success with no env",
			env:  "",
			setup: func() {
				readJsonConfig = func(path string) (*Config, error) {
					return &Config{App: App{Name: "my-service"}}, nil
				}
			},
			want:    &Config{App: App{Name: "my-service"}},
			wantErr: false,
		},
		{
			name: "error due to read JSON config",
			env:  Production,
			setup: func() {
				readJsonConfig = func(path string) (*Config, error) {
					return nil, errors.New("some error")
				}
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Setenv(EnvKey, tt.env)

			tt.setup()

			got, err := LoadConfig()
			if (err != nil) != tt.wantErr {
				t.Errorf("LoadConfig() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("LoadConfig() = %v, want %v", got, tt.want)
			}
		})
	}
}
