package utils

import (
	"os"
	"testing"
)

func TestingRemoveDirectoryTree(t *testing.T, directory string) {
	t.Helper()
	if err := os.RemoveAll(directory); err != nil {
		t.Errorf("Failed to remove the directory tree: %v", err)
	}
}

func TestingRemoveDirectoryTrees(t *testing.T, directories []string) {
	t.Helper()
	for _, directory := range directories {
		TestingRemoveDirectoryTree(t, directory)
	}
}
