package main

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/markuskreukniet/markus-tools/go/utils/test"
)

func getDirectoryPathsAndNamesFromDirectory(directoryFilePath string) ([]string, []string, error) {
	entries, err := os.ReadDir(directoryFilePath)
	if err != nil {
		return nil, nil, err
	}
	var filePaths []string
	var names []string
	for _, entry := range entries {
		if entry.IsDir() {
			names = append(names, entry.Name())
			filePaths = append(filePaths, filepath.Join(directoryFilePath, entry.Name()))
		}
	}
	return filePaths, names, nil
}

func TestFilesToDateRangeDirectory(t *testing.T) {
	// arrange
	// inputAsDelimitedString := `
	// 	,,txt 0.txt,;
	// 	empty,,,;
	// 	directory 1,,txt 1.txt,;
	// 	directory 1,,jpg 1.jpg,;
	// 	directory 2/directory 3,,txt 2 3.txt,;
	// `
	// destinationInputAsDelimitedString := `
	// 	,,txt 0.txt,;
	// 	2020-01-20,,,;
	// 	2020-02-20,,txt 02.txt,;
	// 	2020-03-20,,,;
	// 	2020-04-20 - 2020-04-21,,,;
	// 	2020-05-20 - 2020-05-21,,txt 05.txt,;
	// 	2020-06-20 - 2020-06-21,,,;
	// `
	// inputDirectoriesWithOptionalFile := createDirectoriesWithOptionalFile(inputAsDelimitedString)
	// destinationInputDirectoriesWithOptionalFile := createDirectoriesWithOptionalFile(destinationInputAsDelimitedString)

	// TODO: duplicate naming of files

	testCases := []struct {
		metadata test.TestCaseMetadata
	}{
		{
			metadata: test.TestingCreateTestCaseMetadataWithNameBasicAndWantErrFalse(),
		},
		{
			metadata: test.TestingCreateTestCaseMetadataWithNameEmptyFileSystemNodesAndWantErrFalse(),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.metadata.Name, func(t *testing.T) {
			// arrange and teardown
			// directoryTree := test.TestingCreateTempFileSystemStructureOrGetEmptyString(t, fileSystemPathEndParts)
			// defer test.TestingRemoveDirectoryTree(t, directoryTree)

			// act

			// assert

		})
	}
}
