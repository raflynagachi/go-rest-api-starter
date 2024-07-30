package fs

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type MockFile struct {
	*os.File
	writeErr error
	closeErr error
}

func (m *MockFile) Write(p []byte) (n int, err error) {
	if m.writeErr != nil {
		return 0, m.writeErr
	}
	return len(p), nil
}

func (m *MockFile) Close() error {
	if m.closeErr != nil {
		return m.closeErr
	}
	return nil
}

func TestCreateTempFile(t *testing.T) {
	type args struct {
		content string
		pattern string
	}
	tests := []struct {
		name     string
		args     args
		setup    func()
		wantFile bool
		wantErr  bool
	}{
		{
			name: "success with valid file creation",
			args: args{
				content: "test content",
				pattern: "config-*.json",
			},
			setup: func() {
			},
			wantFile: true,
			wantErr:  false,
		},
		{
			name: "error due to invalid pattern",
			args: args{
				content: "test content",
				pattern: "//invalid//",
			},
			setup:    func() {},
			wantFile: false,
			wantErr:  true,
		},
		// TODO: find a way to mock the os.Write and os.Close
		// {
		// 	name: "error due to invalid write file",
		// 	args: args{
		// 		content: "test content",
		// 		pattern: "config-*.txt",
		// 	},
		// 	setup:    func() {},
		// 	wantFile: false,
		// 	wantErr:  true,
		// },
		// {
		// 	name: "error due to invalid close file",
		// 	args: args{
		// 		content: "test content",
		// 		pattern: "config-*.txt",
		// 	},
		// 	setup: func() {},
		// 	wantFile: false,
		// 	wantErr:  true,
		// },
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()

			gotFile, err := CreateTempFile(tt.args.content, tt.args.pattern)
			if err != nil {
				if !tt.wantErr {
					t.Errorf("CreateTempFile() error = %v, wantErr %v", err, tt.wantErr)
					return
				}
				assert.True(t, tt.wantErr, "Expected an error but got none")
				return
			}

			if gotFile != nil {
				defer os.Remove(gotFile.Name())

				fileContent, err := os.ReadFile(gotFile.Name())
				require.NoError(t, err, "Failed to read file content")
				assert.Equal(t, tt.args.content, string(fileContent), "File content does not match expected")
			}

			if (gotFile != nil) != tt.wantFile {
				t.Errorf("CreateTempFile() = %v, want %v", gotFile, tt.wantFile)
			}
		})
	}
}
