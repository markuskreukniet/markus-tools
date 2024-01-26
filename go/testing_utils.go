package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/markuskreukniet/markus-tools/go/utils/test"
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

// TODO: rename testing to test
// TODO: check if the logic with starting with and without capitals is correct, for example for the functions and vars
// TODO: move this function and other functions to testing_arrange_utils?
func testingCreateTempFileSystemStructureOrGetEmptyString(t *testing.T, fileSystemPathEndParts test.FileSystemPathEndParts) string {
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
