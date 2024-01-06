package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// TODO: rename file?
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

	emptyPathEndParts []string
)

func testingWriteFileContent(t *testing.T, filePath string, content string) {
	if err := os.WriteFile(filePath, []byte(content), 0666); err != nil {
		t.Errorf("Failed to write file content: %v", err)
	}
}

func testingWriteNewlineString(t *testing.T, builder *strings.Builder) {
	_, err := writeNewlineString(builder)
	if err != nil {
		t.Errorf("writeNewlineString failed: %v", err)
	}
}

func testingWriteString(t *testing.T, stringToWrite string, builder *strings.Builder) {
	_, err := builder.WriteString(stringToWrite)
	if err != nil {
		t.Errorf("Failed to write string: %v", err)
	}
}

func testingCreateTempFileSystemStructureOrGetEmptyString(directoryPathEndParts, filePathEndParts []string) (string, error) {
	if len(directoryPathEndParts) == 0 {
		return "", nil
	}
	return testingCreateTempFileSystemStructure(directoryPathEndParts, filePathEndParts)
}

func testingCreateTempFileSystemStructure(directoryPathEndParts, filePathEndParts []string) (string, error) {
	tempDirectory, err := os.MkdirTemp("", "markus-tools go test")
	if err != nil {
		return "", err
	}
	for _, part := range directoryPathEndParts {
		if err := os.MkdirAll(filepath.Join(tempDirectory, part), 0755); err != nil {
			return "", err
		}
	}
	for _, part := range filePathEndParts {
		if err := os.WriteFile(filepath.Join(tempDirectory, part), []byte{}, 0666); err != nil {
			return "", err
		}
	}
	return tempDirectory, nil
}
