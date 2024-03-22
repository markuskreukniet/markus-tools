package main

import (
	"testing"

	"github.com/markuskreukniet/markus-tools/go/utils/test"
)

func TestFilesToDateRangeDirectory(t *testing.T) {
	// arrange
	// source/empty
	// source/directory 1/txt 1.txt
	// source/directory 1/jpg 1.jpg

	// source/directory 1/jpg 1.jpg

	// TODO: duplicate naming of files

	fileSystemPathEndParts := test.FileSystemPathEndParts{
		DirectoryPathEndParts: []string{test.Directory1, test.Directory2WithDirectory3, test.Directory2WithDirectory4},
		FilePathEndParts:      []string{test.TxtFile1, test.TxtFile3, test.TxtFile6, test.JpgFile4},
	}

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
			directoryTree := test.TestingCreateTempFileSystemStructureOrGetEmptyString(t, fileSystemPathEndParts)
			defer test.TestingRemoveDirectoryTree(t, directoryTree)

			// act

			// assert
		})
	}
}
