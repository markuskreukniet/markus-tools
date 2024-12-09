package utils

import (
	"os"
	"testing"
)

func TMustRemoveAll(t *testing.T, filePath string) {
	TMustErr(t, os.RemoveAll(filePath))
}

func RemoveDirectoryTrees(t *testing.T, directoryPaths []string) {
	for _, path := range directoryPaths {
		TMustRemoveAll(t, path)
	}
}
