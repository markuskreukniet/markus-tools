package main

import (
	"testing"

	"github.com/markuskreukniet/markus-tools/go/utils/test"
)

func TestFilesToDateRangeDirectory(t *testing.T) {
	// arrange
	// inputAsDelimitedString := `
	// 	,,txt 0.txt,;
	// 	empty,,nil,;
	// 	directory 1,,txt 1.txt,;
	// 	directory 1,,jpg 1.jpg,;
	// 	directory 2/directory 3,,txt 2 3.txt,;
	// `
	// destinationInputAsDelimitedString := `
	// 	,,txt 0.txt,;
	// 	2020-01-20,,nil,;
	// 	2020-02-20,,txt 02.txt,;
	// 	2020-03-20,,nil,;
	// 	2020-04-20 - 2020-04-21,,nil,;
	// 	2020-05-20 - 2020-05-21,,txt 05.txt,;
	// 	2020-06-20 - 2020-06-21,,nil,;
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
