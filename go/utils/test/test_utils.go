package test

import (
	"strings"
	"testing"
)

func TestingWriteString(t *testing.T, stringToWrite string, builder *strings.Builder) {
	t.Helper()
	_, err := builder.WriteString(stringToWrite)
	if err != nil {
		t.Errorf("Failed to write string: %v", err)
	}
}
