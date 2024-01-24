package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func testingWriteFileContentWithContentAndIndex(t *testing.T, filePath string, index int) string {
	t.Helper()
	writtenContent := fmt.Sprintf("content %d", index)
	testingWriteFileContent(t, filePath, writtenContent)
	return writtenContent
}

func testingWriteFileContent(t *testing.T, filePath string, content string) {
	t.Helper()
	if err := os.WriteFile(filePath, []byte(content), 0666); err != nil {
		t.Errorf("Failed to write file content: %v", err)
	}
}

func testingWriteString(t *testing.T, stringToWrite string, builder *strings.Builder) {
	t.Helper()
	_, err := builder.WriteString(stringToWrite)
	if err != nil {
		t.Errorf("Failed to write string: %v", err)
	}
}

// TODO: move this function and other functions to testing_arrage_utils?
func testingCreateTempFileSystemStructureOrGetEmptyString(t *testing.T, directoryPathEndParts, filePathEndParts []string) string {
	t.Helper()
	if len(directoryPathEndParts) == 0 {
		return ""
	}

	// Create a temporary file system structure.
	tempDirectory, err := os.MkdirTemp("", "markus-tools go test")
	if err != nil {
		// TODO: "Failed to create temp dir" exists. Check also for other duplicate strings
		t.Errorf("Failed to create the temporary directory: %v", err)
	}
	for _, part := range directoryPathEndParts {
		if err := os.MkdirAll(filepath.Join(tempDirectory, part), 0755); err != nil {
			t.Errorf("Failed to create directory in temporary directory: %v", err)
		}
	}
	for _, part := range filePathEndParts {
		if err := os.WriteFile(filepath.Join(tempDirectory, part), []byte{}, 0666); err != nil {
			t.Errorf("Failed to create a file: %v", err)
		}
	}
	return tempDirectory
}
