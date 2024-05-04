package test

import (
	"strings"
	"testing"

	"github.com/markuskreukniet/markus-tools/go/utils"
)

func TestingWriteString(t *testing.T, stringToWrite string, builder *strings.Builder) {
	t.Helper()
	if _, err := builder.WriteString(stringToWrite); err != nil {
		t.Errorf("Failed to write string: %v", err)
	}
}

func TestingWriteTwoNewlineStrings(t *testing.T, builder *strings.Builder) {
	t.Helper()
	if _, err := utils.WriteTwoNewlineStrings(builder); err != nil {
		t.Errorf("WriteTwoNewlineStrings error: %v", err)
	}
}
