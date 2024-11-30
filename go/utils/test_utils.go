package utils

import (
	"strings"
	"testing"
)

func TMust[T any](t *testing.T, v T, err error) T {
	t.Helper()

	if err != nil {
		t.Fatalf("error: %v", err)
	}

	return v
}

func TMustErr(t *testing.T, err error) {
	t.Helper()

	if err != nil {
		t.Fatalf("error: %v", err)
	}
}

func TestingWriteString(t *testing.T, stringToWrite string, builder *strings.Builder) {
	t.Helper()
	if _, err := builder.WriteString(stringToWrite); err != nil {
		t.Errorf("Failed to write string: %v", err)
	}
}

func TestingWriteTwoNewlineStrings(t *testing.T, builder *strings.Builder) {
	t.Helper()
	if _, err := WriteTwoNewlineStrings(builder); err != nil {
		t.Errorf("WriteTwoNewlineStrings error: %v", err)
	}
}
