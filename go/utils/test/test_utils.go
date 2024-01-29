package test

import (
	"os"
	"path/filepath"
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

// TODO: check if the logic with starting with and without capitals is correct, for example for the functions and vars
// TODO: move this function and other functions to testing_arrange_utils?
func TestingCreateTempFileSystemStructureOrGetEmptyString(t *testing.T, fileSystemPathEndParts FileSystemPathEndParts) string {
	t.Helper()
	if len(fileSystemPathEndParts.DirectoryPathEndParts) == 0 {
		return ""
	}

	// Create a temporary file system structure.
	tempDirectory, err := os.MkdirTemp("", "markus-tools go test")
	if err != nil {
		t.Errorf("Failed to create the temporary directory: %v", err)
	}
	for _, part := range fileSystemPathEndParts.DirectoryPathEndParts {
		if err := os.MkdirAll(filepath.Join(tempDirectory, part), 0755); err != nil {
			t.Errorf("Failed to create directory in temporary directory: %v", err)
		}
	}
	for _, part := range fileSystemPathEndParts.FilePathEndParts {
		if err := os.WriteFile(filepath.Join(tempDirectory, part), []byte{}, 0666); err != nil {
			t.Errorf("Failed to create a file: %v", err)
		}
	}
	return tempDirectory
}
