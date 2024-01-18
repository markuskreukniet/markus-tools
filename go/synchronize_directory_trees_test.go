package main

import (
	"bufio"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func testingReadLines(t *testing.T, filePath string) string {
	file, err := os.Open(filePath)
	if err != nil {
		t.Errorf("Failed to open file: %v", err)
	}
	defer file.Close()
	var builder strings.Builder
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		testingWriteString(t, scanner.Text(), &builder)
	}
	err = scanner.Err()
	if err != nil {
		t.Errorf("Failed to read file content: %v", err)
	}
	return builder.String()
}

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

func testingContainsTxtFile4(stringSlice []string) bool {
	for _, item := range stringSlice {
		if item == txtFile4 {
			return true
		}
	}
	return false
}

func TestSynchronizeDirectoryTrees(t *testing.T) {
	// arrange
	sourceDirectoryPathEndParts := []string{directoryEmpty, directory1, directory2WithDirectoryEmpty, directory2WithDirectory3}
	sourceFilePathEndParts := []string{txtFile1, txtFile4, txtFile5}
	destinationDirectoryPathEndParts := []string{directoryEmpty, directory2WithDirectory3}
	destinationFilePathEndParts := []string{txtFile3, txtFile4}

	testCases := []struct {
		Metadata                         TestCaseMetadata
		SourceDirectoryPathEndParts      []string
		SourceFilePathEndParts           []string
		DestinationDirectoryPathEndParts []string
		DestinationFilePathEndParts      []string
		WantSameFilePaths                bool
	}{
		{
			Metadata:                         testingCreateTestCaseMetadataWithNameBasicAndWantErrFalse(),
			SourceDirectoryPathEndParts:      sourceDirectoryPathEndParts,
			SourceFilePathEndParts:           sourceFilePathEndParts,
			DestinationDirectoryPathEndParts: destinationDirectoryPathEndParts,
			DestinationFilePathEndParts:      destinationFilePathEndParts,
			WantSameFilePaths:                true,
		},
		{
			Metadata:                         testingCreateTestCaseMetadataWithWantErrTrue("Empty DestinationPathEndParts"),
			SourceDirectoryPathEndParts:      sourceDirectoryPathEndParts,
			SourceFilePathEndParts:           sourceFilePathEndParts,
			DestinationDirectoryPathEndParts: emptyPathEndParts,
			DestinationFilePathEndParts:      emptyPathEndParts,
			WantSameFilePaths:                false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.Metadata.Name, func(t *testing.T) {
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

			// Some file systems have a resolution of one second, so we must wait a second.
			filePathTxtFile4 := ""
			if sourceDirectory != "" && destinationDirectory != "" && testingContainsTxtFile4(tc.SourceFilePathEndParts) && testingContainsTxtFile4(tc.DestinationFilePathEndParts) {
				filePathTxtFile4 = filepath.Join(destinationDirectory, txtFile4)
				testingWriteFileContentWithContentAndIndex(t, filePathTxtFile4, 1)
				time.Sleep(time.Second)
				testingWriteFileContentWithContentAndIndex(t, filepath.Join(sourceDirectory, txtFile4), 2)
			}

			// act
			err = synchronizeDirectoryTrees(sourceDirectory, destinationDirectory)

			// assert
			testingAssertErrorToWantError(t, err, tc.Metadata.WantErr)
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
			// TODO: "content 2" should have returned testingWriteFileContentWithContentAndIndex?
			if filePathTxtFile4 != "" && testingReadLines(t, filePathTxtFile4) != "content 2" {
				// TODO:
				t.Fatalf("--- fout ---")
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
		Metadata       TestCaseMetadata
		InputBasePath  string
		InputFullPath  string
		OutputBasePath string
		Want           string
	}{
		{
			Metadata:       testingCreateTestCaseMetadataWithNameBasicAndWantErrFalse(),
			InputBasePath:  inputBasePath,
			InputFullPath:  inputFullPath,
			OutputBasePath: outputBasePath,
			Want:           filepath.FromSlash(joinedOutputBasePathWithRelativeInputPath),
		},
		{
			Metadata:       testingCreateTestCaseMetadataWithWantErrTrue("Empty InputBasePath"),
			InputBasePath:  "",
			InputFullPath:  inputFullPath,
			OutputBasePath: outputBasePath,
			Want:           "",
		},
		{
			Metadata:       testingCreateTestCaseMetadataWithWantErrTrue("Empty InputFullPath"),
			InputBasePath:  inputBasePath,
			InputFullPath:  "",
			OutputBasePath: outputBasePath,
			Want:           "",
		},
		{
			Metadata:       testingCreateTestCaseMetadata("Equivalent Input Paths", false),
			InputBasePath:  inputBasePath,
			InputFullPath:  inputBasePath,
			OutputBasePath: outputBasePath,
			Want:           filepath.FromSlash(outputBasePath),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.Metadata.Name, func(t *testing.T) {
			result, err := joinOutputBasePathWithRelativeInputPath(tc.InputBasePath, tc.InputFullPath, tc.OutputBasePath)
			testingAssertErrorToWantError(t, err, tc.Metadata.WantErr)
			if err == nil && result != tc.Want {
				t.Fatalf("want: %s, got %s", tc.Want, result)
			}
		})
	}
}
