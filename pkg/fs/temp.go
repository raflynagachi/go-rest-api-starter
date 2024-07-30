package fs

import (
	"os"

	"github.com/pkg/errors"
)

// CreateTempFile create temporary file with the given content.
// It is the caller's responsibility to remove the file when it is no longer needed.
func CreateTempFile(content string, pattern string) (file *os.File, err error) {
	file, err = os.CreateTemp("", pattern)
	if err != nil {
		return nil, errors.Errorf("Failed to create temporary file: %v", err)
	}

	if _, err := file.Write([]byte(content)); err != nil {
		return nil, errors.Errorf("Failed to write to temporary file: %v", err)
	}

	if err := file.Close(); err != nil {
		return nil, errors.Errorf("Failed to close temporary file: %v", err)
	}

	return file, nil
}
