package main

import (
	"os"
	"testing"
)

func testingRemoveDirectoryTree(t *testing.T, directory string) {
	t.Helper()
	if err := os.RemoveAll(directory); err != nil {
		t.Errorf("Failed to remove the directory tree: %v", err)
	}
}
