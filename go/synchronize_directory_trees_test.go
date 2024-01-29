package main

import (
	"bufio"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/markuskreukniet/markus-tools/go/utils/test"
)

func testingReadLines(t *testing.T, filePath string) string {
	t.Helper()
	file, err := os.Open(filePath)
	if err != nil {
		t.Errorf("Failed to open file: %v", err)
	}
	defer file.Close()
	var builder strings.Builder
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		test.TestingWriteString(t, scanner.Text(), &builder)
	}
	err = scanner.Err()
	if err != nil {
		t.Errorf("Failed to read file content: %v", err)
	}
	return builder.String()
}

func testingHaveDirectoryTreesSameFilePathsOrGetFalse(t *testing.T, sourceDirectory, destinationDirectory string) bool {
	t.Helper()
	if sourceDirectory == "" || destinationDirectory == "" {
		return false
	}

	// Do the directory trees have the same file paths?
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
		t.Errorf("Failed to walk through directory: %v", err)
	}
	return haveSameFilePaths
}

func testingContainsTxtFile4(stringSlice []string) bool {
	for _, item := range stringSlice {
		if item == test.TxtFile4 {
			return true
		}
	}
	return false
}

func TestSynchronizeDirectoryTrees(t *testing.T) {
	// arrange
	sourceFileSystemPathEndParts := test.FileSystemPathEndParts{
		DirectoryPathEndParts: []string{test.DirectoryEmpty, test.Directory1, test.Directory2WithDirectoryEmpty, test.Directory2WithDirectory3},
		FilePathEndParts:      []string{test.TxtFile1, test.TxtFile4, test.TxtFile5},
	}
	destinationFileSystemPathEndParts := test.FileSystemPathEndParts{
		DirectoryPathEndParts: []string{test.DirectoryEmpty, test.Directory2WithDirectory3},
		FilePathEndParts:      []string{test.TxtFile3, test.TxtFile4},
	}

	testCases := []struct {
		Metadata                          test.TestCaseMetadata
		SourceFileSystemPathEndParts      test.FileSystemPathEndParts
		DestinationFileSystemPathEndParts test.FileSystemPathEndParts
		WantSameFilePaths                 bool
	}{
		{
			Metadata:                          test.TestingCreateTestCaseMetadataWithNameBasicAndWantErrFalse(),
			SourceFileSystemPathEndParts:      sourceFileSystemPathEndParts,
			DestinationFileSystemPathEndParts: destinationFileSystemPathEndParts,
			WantSameFilePaths:                 true,
		},
		{
			Metadata:                          test.TestingCreateTestCaseMetadataWithWantErrTrue("Empty DestinationPathEndParts"),
			SourceFileSystemPathEndParts:      sourceFileSystemPathEndParts,
			DestinationFileSystemPathEndParts: test.EmptyFileSystemPathEndParts,
			WantSameFilePaths:                 false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.Metadata.Name, func(t *testing.T) {
			// arrange and teardown
			sourceDirectory := test.TestingCreateTempFileSystemStructureOrGetEmptyString(t, tc.SourceFileSystemPathEndParts)
			defer test.TestingRemoveDirectoryTree(t, sourceDirectory)
			destinationDirectory := test.TestingCreateTempFileSystemStructureOrGetEmptyString(t, tc.DestinationFileSystemPathEndParts)
			defer test.TestingRemoveDirectoryTree(t, destinationDirectory)

			// Some file systems have a resolution of one second, so we must wait a second.
			filePathTxtFile4 := ""
			writtenContent := ""
			if sourceDirectory != "" && destinationDirectory != "" && testingContainsTxtFile4(tc.SourceFileSystemPathEndParts.FilePathEndParts) && testingContainsTxtFile4(tc.DestinationFileSystemPathEndParts.FilePathEndParts) {
				filePathTxtFile4 = filepath.Join(destinationDirectory, test.TxtFile4)
				test.TestingWriteFileContentWithContentAndIndex(t, filePathTxtFile4, 1)
				time.Sleep(time.Second)
				writtenContent = test.TestingWriteFileContentWithContentAndIndex(t, filepath.Join(sourceDirectory, test.TxtFile4), 2)
			}

			// act
			err := synchronizeDirectoryTrees(sourceDirectory, destinationDirectory)

			// assert
			test.TestingAssertErrorToWantError(t, err, tc.Metadata.WantErr)
			haveSameFilePaths := testingHaveDirectoryTreesSameFilePathsOrGetFalse(t, sourceDirectory, destinationDirectory)
			if tc.WantSameFilePaths && !haveSameFilePaths {
				t.Errorf("The source and destination directory trees do not have the same file paths.")
			}
			haveSameFilePaths = testingHaveDirectoryTreesSameFilePathsOrGetFalse(t, destinationDirectory, sourceDirectory)
			if tc.WantSameFilePaths && !haveSameFilePaths {
				t.Errorf("The destination and source directory trees do not have the same file paths.")
			}
			if filePathTxtFile4 != "" {
				test.TestingAssertEqualStrings(t, testingReadLines(t, filePathTxtFile4), writtenContent)
			}
		})
	}
}

// TODO: use vars from arrange utils?
func TestJoinOutputBasePathWithRelativeInputPath(t *testing.T) {
	const inputBasePath string = "/home/user/source"
	const inputFullPath string = "/home/user/source/directory/file.txt"
	const outputBasePath string = "/home/user/destination"
	const joinedOutputBasePathWithRelativeInputPath string = "/home/user/destination/directory/file.txt"

	testCases := []struct {
		Metadata       test.TestCaseMetadata
		InputBasePath  string
		InputFullPath  string
		OutputBasePath string
		Want           string
	}{
		{
			Metadata:       test.TestingCreateTestCaseMetadataWithNameBasicAndWantErrFalse(),
			InputBasePath:  inputBasePath,
			InputFullPath:  inputFullPath,
			OutputBasePath: outputBasePath,
			Want:           filepath.FromSlash(joinedOutputBasePathWithRelativeInputPath),
		},
		{
			Metadata:       test.TestingCreateTestCaseMetadataWithWantErrTrue("Empty InputBasePath"),
			InputBasePath:  "",
			InputFullPath:  inputFullPath,
			OutputBasePath: outputBasePath,
			Want:           "",
		},
		{
			Metadata:       test.TestingCreateTestCaseMetadataWithWantErrTrue("Empty InputFullPath"),
			InputBasePath:  inputBasePath,
			InputFullPath:  "",
			OutputBasePath: outputBasePath,
			Want:           "",
		},
		{
			Metadata:       test.TestingCreateTestCaseMetadata("Equivalent Input Paths", false),
			InputBasePath:  inputBasePath,
			InputFullPath:  inputBasePath,
			OutputBasePath: outputBasePath,
			Want:           filepath.FromSlash(outputBasePath),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.Metadata.Name, func(t *testing.T) {
			result, err := joinOutputBasePathWithRelativeInputPath(tc.InputBasePath, tc.InputFullPath, tc.OutputBasePath)
			test.TestingAssertErrorToWantError(t, err, tc.Metadata.WantErr)
			if err == nil {
				test.TestingAssertEqualStrings(t, tc.Want, result)
			}
		})
	}
}
