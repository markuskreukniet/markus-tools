package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// TODO: should move to arrange utils?
// TODO: there are no files in the root (temp dir)
var (
	directoryEmpty               = "directory empty"
	directory1                   = "directory 1"
	directory2                   = "directory 2"
	directory2WithDirectoryEmpty = filepath.Join(directory2, directoryEmpty)
	directory2WithDirectory3     = filepath.Join(directory2, "directory 3")
	directory2WithDirectory4     = filepath.Join(directory2, "directory 4")

	txtFile1 = filepath.Join(directory1, "file 1.txt")
	txtFile2 = filepath.Join(directory1, "file 2.txt")
	txtFile3 = filepath.Join(directory2WithDirectory3, "file 3.txt")
	txtFile4 = filepath.Join(directory2WithDirectory3, "file 4.txt")
	txtFile5 = filepath.Join(directory2WithDirectory3, "file 5.txt")
	txtFile6 = filepath.Join(directory2WithDirectory4, "file 6.txt")

	jpgFile4 = filepath.Join(directory1, "file 4.jpg")

	txtFileNonExistent1 = "non existent 1.txt"

	emptyPathEndParts []string
)

// TODO: move to other util file, or don't return error and use t.Errorf?
func writeNewlineString(builder *strings.Builder) (int, error) {
	bytesWritten, err := builder.WriteString("\n")
	if err != nil {
		return bytesWritten, err
	}
	return bytesWritten, nil
}

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

func testingWriteNewlineString(t *testing.T, builder *strings.Builder) {
	t.Helper()
	_, err := writeNewlineString(builder)
	if err != nil {
		t.Errorf("writeNewlineString failed: %v", err)
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
	return testingCreateTempFileSystemStructure(t, directoryPathEndParts, filePathEndParts)
}

// TODO: useless function?
func testingCreateTempFileSystemStructure(t *testing.T, directoryPathEndParts, filePathEndParts []string) string {
	tempDirectory, err := os.MkdirTemp("", "markus-tools go test")
	if err != nil {
		// TODO: "Failed to create temp dir" exists. Check also for other duplicate strings
		t.Fatalf("Failed to create the temporary directory: %v", err)
	}
	for _, part := range directoryPathEndParts {
		if err := os.MkdirAll(filepath.Join(tempDirectory, part), 0755); err != nil {
			t.Fatalf("Failed to create directory in temporary directory: %v", err)
		}
	}
	for _, part := range filePathEndParts {
		if err := os.WriteFile(filepath.Join(tempDirectory, part), []byte{}, 0666); err != nil {
			t.Fatalf("Failed to create a file: %v", err)
		}
	}
	return tempDirectory
}
