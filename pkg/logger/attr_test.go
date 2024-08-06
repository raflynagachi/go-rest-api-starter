package logger

import (
	"fmt"
	"log/slog"
	"testing"
	"time"
)

func createTestRecord(level slog.Level, message string, attrs ...slog.Attr) slog.Record {
	rec := slog.NewRecord(time.Now(), level, message, 0)
	for _, attr := range attrs {
		rec.Add(attr)
	}
	return rec
}

func TestFormatAttributes(t *testing.T) {
	rec := slog.NewRecord(time.Now(), LevelInfo, "test message", 0)
	rec.Add(slog.String("key", "value"))

	expected := map[string]interface{}{
		"key": "value",
	}

	result := formatAttributes(rec)

	for k, v := range expected {
		if result[k] != v {
			t.Errorf("formatAttributes() = %v, want %v", result, expected)
		}
	}
}

func TestGroupAttr(t *testing.T) {
	attr1 := slog.String("key1", "value1")
	attr2 := slog.Int("key2", 42)
	attr3 := slog.Float64("key3", 3.14)

	group := GroupAttr("groupKey", attr1, attr2, attr3)
	expected := slog.Group("groupKey", attr1, attr2, attr3)

	if !equalAttrs(group, expected) {
		t.Errorf("GroupAttr() = %v, want %v", group, expected)
	}
}

func equalAttrs(a, b slog.Attr) bool {
	if a.Key != b.Key {
		return false
	}

	groupA, okA := a.Value.Any().([]slog.Attr)
	groupB, okB := b.Value.Any().([]slog.Attr)
	if !okA || !okB {
		return false
	}
	if len(groupA) != len(groupB) {
		return false
	}

	for i := range groupA {
		if !groupA[i].Equal(groupB[i]) {
			return false
		}
	}
	return true
}

func TestStringAttr(t *testing.T) {
	attr := StringAttr("stringKey", "test")
	expected := slog.String("stringKey", "test")

	if !attr.Equal(expected) {
		t.Errorf("StringAttr() = %v, want %v", attr, expected)
	}
}

func TestFloat64Attr(t *testing.T) {
	attr := Float64Attr("float64Key", 1.23)
	expected := slog.Float64("float64Key", 1.23)

	if !attr.Equal(expected) {
		t.Errorf("Float64Attr() = %v, want %v", attr, expected)
	}
}

func TestInt64Attr(t *testing.T) {
	attr := Int64Attr("int64Key", 1234567890)
	expected := slog.Int64("int64Key", 1234567890)

	if !attr.Equal(expected) {
		t.Errorf("Int64Attr() = %v, want %v", attr, expected)
	}
}

func TestUint64Attr(t *testing.T) {
	attr := Uint64Attr("int64Key", 1234567890)
	expected := slog.Uint64("int64Key", 1234567890)

	if !attr.Equal(expected) {
		t.Errorf("Int64Attr() = %v, want %v", attr, expected)
	}
}

func TestTimeAttr(t *testing.T) {
	now := time.Now()
	attr := TimeAttr("timeKey", now)
	expected := slog.String("timeKey", now.Format(time.RFC3339))

	if !attr.Equal(expected) {
		t.Errorf("TimeAttr() = %v, want %v", attr, expected)
	}
}

func TestErrAttr(t *testing.T) {
	tests := []struct {
		name     string
		err      error
		expected string
	}{
		{name: "Non-nil error", err: fmt.Errorf("test error"), expected: "test error"},
		{name: "Nil error", err: nil, expected: "no error"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			attr := ErrAttr(tt.err)
			expected := slog.String("error", tt.expected)

			if !attr.Equal(expected) {
				t.Errorf("ErrAttr() = %v, want %v", attr, expected)
			}
		})
	}
}
