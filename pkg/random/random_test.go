package random

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRandomInt(t *testing.T) {
	tests := []struct {
		name string
		min  int
		max  int
	}{
		{"normal range", 1, 10},
		{"min equals max", 5, 5},
		{"negative range", -10, -1},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := RandomInt(tt.min, tt.max)

			assert.GreaterOrEqual(t, result, tt.min, "Result should be >= min")
			assert.LessOrEqual(t, result, tt.max, "Result should be <= max")
		})
	}
}

func TestRandomID(t *testing.T) {
	tests := []struct {
		name string
	}{
		{"random ID"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := RandomID()

			assert.Greater(t, result, int64(0), "ID should be > 0")
			assert.LessOrEqual(t, result, int64(1000000000), "ID should be <= 1,000,000,000")
		})
	}
}

func TestRandomFloat(t *testing.T) {
	tests := []struct {
		name string
		min  int
		max  int
	}{
		{"normal range", 1, 10},
		{"min equals max", 5, 5},
		{"negative range", -10, -1},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := RandomFloat(tt.min, tt.max)

			assert.GreaterOrEqual(t, result, float64(tt.min), "Result should be >= min")
			assert.LessOrEqual(t, result, float64(tt.max), "Result should be <= max")
		})
	}
}

func TestRandomString(t *testing.T) {
	tests := []struct {
		name string
		len  int
	}{
		{"short string", 5},
		{"medium string", 10},
		{"long string", 20},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := RandomString(tt.len)

			assert.Equal(t, tt.len, len(result), "String length should be equal to requested length")
			assert.True(t, strings.IndexFunc(result, func(r rune) bool {
				return !strings.ContainsRune(alphabet, r)
			}) == -1, "String should only contain characters from the alphabet")
		})
	}
}

func TestRandomEmail(t *testing.T) {
	tests := []struct {
		name string
	}{
		{"random email"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := RandomEmail()

			assert.True(t, strings.HasSuffix(result, "@mail.com"), "Email should end with @mail.com")
			assert.Len(t, result[:strings.Index(result, "@")], 6, "Email prefix should be 6 characters long")
		})
	}
}
