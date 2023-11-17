package main

import (
	"os"
	"path/filepath"
	"strconv"
	"testing"
)

// TODO: Might work
func TestGetDuplicateFilesAsNewlineSeparatedString(t *testing.T) {
	// Step 1: Test Setup
	tempDir, err := os.MkdirTemp("", "test")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tempDir) // Clean up

	// Create a set of test files, some of which are duplicates
	fileContents := []string{"content1", "content2", "content1"} // 'content1' is duplicated
	filePaths := make([]string, len(fileContents))
	for i, content := range fileContents {
		filePath := filepath.Join(tempDir, "file"+strconv.Itoa(i))
		if err := os.WriteFile(filePath, []byte(content), 0666); err != nil {
			t.Fatalf("Failed to create test file: %v", err)
		}
		filePaths[i] = filePath
	}

	// Step 2: Test Execution
	// Convert filePaths to FileSystemNodes
	nodes := make([]FileSystemNode, len(filePaths))
	for i, path := range filePaths {
		nodes[i] = FileSystemNode{Path: path, IsDirectory: false}
	}

	result, err := getDuplicateFilesAsNewlineSeparatedString(nodes)
	if err != nil {
		t.Fatalf("Function returned an error: %v", err)
	}

	// Step 3: Verification
	expectedResult := filePaths[0] + "\n" + filePaths[2]
	if result != expectedResult {
		t.Errorf("Expected %q, got %q", expectedResult, result)
	}

	// Step 4: Test Teardown is handled by defer statement
}
