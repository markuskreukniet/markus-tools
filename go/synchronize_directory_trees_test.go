package main

import (
	"io/fs"
	"os"
	"path/filepath"
	"testing"
)

// WIP
func TestSynchronizeDirectoryTrees(t *testing.T) {
	testCases := []struct {
		name          string
		sourceStruct  map[string]fs.FileMode
		destStruct    map[string]fs.FileMode
		expectedError bool
	}{
		{
			name: "Simple sync",
			sourceStruct: map[string]fs.FileMode{
				"file1.txt": 0666,
				"dir":       0755 | fs.ModeDir,
			},
			destStruct:    map[string]fs.FileMode{},
			expectedError: false,
		},
		{
			name: "Simple sync 2",
			sourceStruct: map[string]fs.FileMode{
				"file1.txt": 0666,
				"dir":       0755 | fs.ModeDir,
			},
			destStruct: map[string]fs.FileMode{
				"file2.txt": 0666,
				"dir2":      0755 | fs.ModeDir,
			},
			expectedError: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			sourceDir, err := os.MkdirTemp("", "source")
			if err != nil {
				t.Fatalf("Failed to create temporary source directory: %v", err)
			}
			destDir, err := os.MkdirTemp("", "dest")
			if err != nil {
				t.Fatalf("Failed to create temporary destination directory: %v", err)
			}

			if err := createTestDirectoryStructure(sourceDir, tc.sourceStruct); err != nil {
				t.Fatalf("Failed to create source directory structure: %v", err)
			}
			if err := createTestDirectoryStructure(destDir, tc.destStruct); err != nil {
				t.Fatalf("Failed to create destination directory structure: %v", err)
			}

			err = synchronizeDirectoryTrees(sourceDir, destDir)

			if (err != nil) != tc.expectedError {
				t.Errorf("synchronizeDirectoryTrees() error = %v, expectedError %v", err, tc.expectedError)
			}

			// Check if the destination directory has the expected files and directories
			err = filepath.Walk(destDir, func(path string, info os.FileInfo, err error) error {
				if err != nil {
					return err
				}

				relPath, err := filepath.Rel(destDir, path)
				if err != nil {
					return err
				}

				// Skip the root of the destination directory
				if relPath == "." {
					return nil
				}

				srcPath := filepath.Join(sourceDir, relPath)
				if _, err := os.Stat(srcPath); os.IsNotExist(err) {
					t.Errorf("Extra file or directory in destination: %s", relPath)
				}

				return nil
			})
			if err != nil {
				t.Errorf("Error walking through destination directory: %v", err)
			}

			// Optionally, check file contents and other properties
			// ...

			if err := os.RemoveAll(sourceDir); err != nil {
				t.Fatalf("Failed to remove source directory: %v", err)
			}
			if err := os.RemoveAll(destDir); err != nil {
				t.Fatalf("Failed to remove destination directory: %v", err)
			}
		})
	}
}
