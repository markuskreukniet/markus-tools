package main

import (
	"os"
	"path/filepath"
	"testing"
)

func testingHaveDirectoryTreesSameFilePathsOrGetFalse(sourceDirectory, destinationDirectory string) (bool, error) {
	if sourceDirectory == "" || destinationDirectory == "" {
		return false, nil
	}
	return testingHaveDirectoryTreesSameFilePaths(sourceDirectory, destinationDirectory)
}

func testingHaveDirectoryTreesSameFilePaths(sourceDirectory, destinationDirectory string) (bool, error) {
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

func TestSynchronizeDirectoryTrees(t *testing.T) {
	// arrange
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
			sourceDirectory, err := testingCreateTempFileSystemStructureOrGetEmptyString(tc.SourceDirectoryPathEndParts, tc.SourceFilePathEndParts)
			if err != nil {
				t.Fatalf("Failed to create the temporary source directory: %v", err)
			}
			defer func() {
				if err := os.RemoveAll(sourceDirectory); err != nil {
					t.Errorf("Failed to remove the temporary source directory: %v", err)
				}
			}()
			destinationDirectory, err := testingCreateTempFileSystemStructureOrGetEmptyString(tc.DestinationDirectoryPathEndParts, tc.DestinationFilePathEndParts)
			if err != nil {
				t.Fatalf("Failed to create the temporary destination directory: %v", err)
			}
			defer func() {
				if err := os.RemoveAll(destinationDirectory); err != nil {
					t.Errorf("Failed to remove the temporary destination directory: %v", err)
				}
			}()

			// act
			err = synchronizeDirectoryTrees(sourceDirectory, destinationDirectory)

			// assert
			if (err != nil) != tc.WantErr {
				t.Fatalf("want error: %v, got %v", tc.WantErr, err)
			}
			haveSameFilePaths, err := testingHaveDirectoryTreesSameFilePathsOrGetFalse(sourceDirectory, destinationDirectory)
			if err != nil {
				t.Fatalf("Failed to check if the source and destination directory trees have the same file paths: %v", err)
			}
			if tc.WantSameFilePaths && !haveSameFilePaths {
				t.Fatalf("The source and destination directory trees do not have the same file paths.")
			}
			haveSameFilePaths, err = testingHaveDirectoryTreesSameFilePathsOrGetFalse(destinationDirectory, sourceDirectory)
			if err != nil {
				t.Fatalf("Failed to check if the destination and source directory trees have the same file paths: %v", err)
			}
			if tc.WantSameFilePaths && !haveSameFilePaths {
				t.Fatalf("The destination and source directory trees do not have the same file paths.")
			}
		})
	}
}

func TestJoinOutputBasePathWithRelativeInputPath(t *testing.T) {
	const inputBasePath string = "/home/user/source"
	const inputFullPath string = "/home/user/source/directory/file.txt"
	const outputBasePath string = "/home/user/destination"
	const joinedOutputBasePathWithRelativeInputPath string = "/home/user/destination/directory/file.txt"

	testCases := []struct {
		Name           string
		InputBasePath  string
		InputFullPath  string
		OutputBasePath string
		Want           string
		WantErr        bool
	}{
		{
			Name:           "Basic",
			InputBasePath:  inputBasePath,
			InputFullPath:  inputFullPath,
			OutputBasePath: outputBasePath,
			Want:           filepath.FromSlash(joinedOutputBasePathWithRelativeInputPath),
			WantErr:        false,
		},
		{
			Name:           "Empty inputBasePath",
			InputBasePath:  "",
			InputFullPath:  inputFullPath,
			OutputBasePath: outputBasePath,
			Want:           "",
			WantErr:        true,
		},
		{
			Name:           "Empty inputFullPath",
			InputBasePath:  inputBasePath,
			InputFullPath:  "",
			OutputBasePath: outputBasePath,
			Want:           "",
			WantErr:        true,
		},
		{
			Name:           "Equivalent Input Paths",
			InputBasePath:  inputBasePath,
			InputFullPath:  inputBasePath,
			OutputBasePath: outputBasePath,
			Want:           filepath.FromSlash(outputBasePath),
			WantErr:        false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			result, err := joinOutputBasePathWithRelativeInputPath(tc.InputBasePath, tc.InputFullPath, tc.OutputBasePath)
			if (err != nil) != tc.WantErr {
				t.Fatalf("want error: %v, got %v", tc.WantErr, err)
			}
			if err == nil && result != tc.Want {
				t.Fatalf("want: %s, got %s", tc.Want, result)
			}
		})
	}
}
