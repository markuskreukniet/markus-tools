package main

import (
	"os"
	"path/filepath"
	"testing"
)

func createTempFileSystemStructure(directoryPathEndParts []string, filePathEndParts []string) (string, error) {
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

func TestSynchronizeDirectoryTrees(t *testing.T) {
	// arrange
	directoryEmpty := "directory empty"
	directory1 := "directory 1"
	directory2 := "directory 2"
	directory2WithDirectoryEmpty := filepath.Join(directory2, directoryEmpty)
	directory2WithDirectory3 := filepath.Join(directory2, "directory 3")

	txtFile1 := filepath.Join(directory1, "file 1.txt")
	txtFile2 := filepath.Join(directory1, "file 2.txt")
	txtFile3 := filepath.Join(directory2WithDirectory3, "file 3.txt")
	txtFile4 := filepath.Join(directory2WithDirectory3, "file 4.txt")

	sourceDirectoryPathEndParts := []string{directoryEmpty, directory1, directory2WithDirectoryEmpty, directory2WithDirectory3}
	sourceFilePathEndParts := []string{txtFile1, txtFile2, txtFile3}

	destinationDirectoryPathEndParts := []string{directoryEmpty, directory2WithDirectory3}
	destinationFilePathEndParts := []string{txtFile3, txtFile4}

	testCases := []struct {
		Name                             string
		SourceDirectoryPathEndParts      []string
		SourceFilePathEndParts           []string
		DestinationDirectoryPathEndParts []string
		DestinationFilePathEndParts      []string
		WantErr                          bool
	}{
		{
			Name:                             "Basic",
			SourceDirectoryPathEndParts:      sourceDirectoryPathEndParts,
			SourceFilePathEndParts:           sourceFilePathEndParts,
			DestinationDirectoryPathEndParts: destinationDirectoryPathEndParts,
			DestinationFilePathEndParts:      destinationFilePathEndParts,
			WantErr:                          false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			// arrange
			sourceDirectory, err := createTempFileSystemStructure(tc.SourceDirectoryPathEndParts, tc.SourceFilePathEndParts)
			if err != nil {
				t.Fatalf("Failed to create temporary source directory: %v", err)
			}
			destinationDirectory, err := createTempFileSystemStructure(tc.DestinationDirectoryPathEndParts, tc.DestinationFilePathEndParts)
			if err != nil {
				t.Fatalf("Failed to create temporary destination directory: %v", err)
			}

			// act
			err = synchronizeDirectoryTrees(sourceDirectory, destinationDirectory)

			// assert
			if (err != nil) != tc.WantErr {
				t.Fatalf("want error: %v, got %v", tc.WantErr, err)
			}

			// TODO: this walk code with error should a function that is also tested?
			// TODO: Also, now is only checked if the source has the same files/folders as destination, but not the other way around. Maybe a file or folder did not get copied
			err = filepath.Walk(destinationDirectory, func(path string, info os.FileInfo, err error) error {
				if err != nil {
					return err
				}
				relPath, err := filepath.Rel(destinationDirectory, path)
				if err != nil {
					return err
				}
				if relPath != "." {
					sourcePath := filepath.Join(sourceDirectory, relPath)
					if _, err := os.Stat(sourcePath); os.IsNotExist(err) {
						t.Errorf("Extra file or directory in destination: %s", relPath)
					}
				}
				return nil
			})
			if err != nil {
				t.Errorf("Error walking through destination directory: %v", err)
			}

			// tear down
			if err := os.RemoveAll(sourceDirectory); err != nil {
				t.Fatalf("Failed to remove source directory: %v", err)
			}
			if err := os.RemoveAll(destinationDirectory); err != nil {
				t.Fatalf("Failed to remove destination directory: %v", err)
			}
		})
	}
}
