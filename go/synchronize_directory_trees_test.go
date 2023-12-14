package main

import (
	"os"
	"path/filepath"
	"testing"
)

func haveDirectoryTreesSameFilePathsOrGetFalse(sourceDirectory, destinationDirectory string) (bool, error) {
	if sourceDirectory == "" || destinationDirectory == "" {
		return false, nil
	}
	return haveDirectoryTreesSameFilePaths(sourceDirectory, destinationDirectory)
}

func haveDirectoryTreesSameFilePaths(sourceDirectory, destinationDirectory string) (bool, error) {
	haveSameFilePaths := true
	err := filepath.Walk(sourceDirectory, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		relativePath, err := filepath.Rel(sourceDirectory, path)
		if err != nil {
			return err
		}
		destinationPath := filepath.Join(destinationDirectory, relativePath)
		if _, err := os.Stat(destinationPath); os.IsNotExist(err) {
			haveSameFilePaths = false
			return err
		}
		return nil
	})
	if err != nil {
		return false, err
	}
	return haveSameFilePaths, nil
}

func createTempFileSystemStructureOrGetEmptyString(directoryPathEndParts, filePathEndParts []string) (string, error) {
	if len(directoryPathEndParts) == 0 {
		return "", nil
	}
	return createTempFileSystemStructure(directoryPathEndParts, filePathEndParts)
}

func createTempFileSystemStructure(directoryPathEndParts, filePathEndParts []string) (string, error) {
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

	var emptyPathEndParts []string

	testCases := []struct {
		Name                             string
		SourceDirectoryPathEndParts      []string
		SourceFilePathEndParts           []string
		DestinationDirectoryPathEndParts []string
		DestinationFilePathEndParts      []string
		WantSameFilePaths                bool
		WantErr                          bool
	}{
		{
			Name:                             "Basic",
			SourceDirectoryPathEndParts:      sourceDirectoryPathEndParts,
			SourceFilePathEndParts:           sourceFilePathEndParts,
			DestinationDirectoryPathEndParts: destinationDirectoryPathEndParts,
			DestinationFilePathEndParts:      destinationFilePathEndParts,
			WantSameFilePaths:                true,
			WantErr:                          false,
		},
		{
			Name:                             "Empty DestinationPathEndParts",
			SourceDirectoryPathEndParts:      sourceDirectoryPathEndParts,
			SourceFilePathEndParts:           sourceFilePathEndParts,
			DestinationDirectoryPathEndParts: emptyPathEndParts,
			DestinationFilePathEndParts:      emptyPathEndParts,
			WantSameFilePaths:                false,
			WantErr:                          true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			// arrange and tear down
			sourceDirectory, err := createTempFileSystemStructureOrGetEmptyString(tc.SourceDirectoryPathEndParts, tc.SourceFilePathEndParts)
			if err != nil {
				t.Fatalf("Failed to create temporary source directory: %v", err)
			}
			defer func() {
				if err := os.RemoveAll(sourceDirectory); err != nil {
					t.Errorf("Failed to remove temporary source directory: %v", err)
				}
			}()
			destinationDirectory, err := createTempFileSystemStructureOrGetEmptyString(tc.DestinationDirectoryPathEndParts, tc.DestinationFilePathEndParts)
			if err != nil {
				t.Fatalf("Failed to create temporary destination directory: %v", err)
			}
			defer func() {
				if err := os.RemoveAll(destinationDirectory); err != nil {
					t.Errorf("Failed to remove temporary destination directory: %v", err)
				}
			}()

			// act
			err = synchronizeDirectoryTrees(sourceDirectory, destinationDirectory)

			// assert
			if (err != nil) != tc.WantErr {
				t.Fatalf("want error: %v, got %v", tc.WantErr, err)
			}
			haveSameFilePaths, err := haveDirectoryTreesSameFilePathsOrGetFalse(sourceDirectory, destinationDirectory)
			if err != nil {
				t.Fatalf("Failed to check if the source and destination directory trees have the same file paths: %v", err)
			}
			if tc.WantSameFilePaths && !haveSameFilePaths {
				t.Fatalf("The source and destination directory trees do not have the same file paths.")
			}
			haveSameFilePaths, err = haveDirectoryTreesSameFilePathsOrGetFalse(destinationDirectory, sourceDirectory)
			if err != nil {
				t.Fatalf("Failed to check if the destination and source directory trees have the same file paths: %v", err)
			}
			if tc.WantSameFilePaths && !haveSameFilePaths {
				t.Fatalf("The destination and source directory trees do not have the same file paths.")
			}
		})
	}
}
