package main

import (
	"path/filepath"
	"strings"
	"testing"

	"github.com/markuskreukniet/markus-tools/go/utils"
	"github.com/markuskreukniet/markus-tools/go/utils/test"
)

func testingWriteNewlineString(t *testing.T, builder *strings.Builder) {
	t.Helper()
	_, err := utils.WriteNewlineString(builder)
	if err != nil {
		t.Errorf("writeNewlineString error: %v", err)
	}
}

func TestGetDuplicateFilesAsNewlineSeparatedString(t *testing.T) {
	// arrange
	fileSystemPathEndParts := test.FileSystemPathEndParts{
		DirectoryPathEndParts: []string{test.DirectoryEmpty, test.Directory1, test.Directory2WithDirectoryEmpty, test.Directory2WithDirectory3, test.Directory2WithDirectory4},
		FilePathEndParts:      []string{test.TxtFile1, test.TxtFile2, test.TxtFile3, test.TxtFile4, test.TxtFile5, test.TxtFile6},
	}
	duplicateFilePathEndPartGroups := [][]string{{test.TxtFile2, test.TxtFile3}, {test.TxtFile4, test.TxtFile5, test.TxtFile6}}
	var emptyPathEndPartGroups [][]string

	testCases := []struct {
		Metadata                       test.TestCaseMetadata
		FileSystemPathEndParts         test.FileSystemPathEndParts
		DuplicateFilePathEndPartGroups [][]string
	}{
		{
			Metadata:                       test.TestingCreateTestCaseMetadataWithNameBasicAndWantErrFalse(),
			FileSystemPathEndParts:         fileSystemPathEndParts,
			DuplicateFilePathEndPartGroups: duplicateFilePathEndPartGroups,
		},
		{
			Metadata:                       test.TestingCreateTestCaseMetadataWithNameEmptyFileSystemNodesAndWantErrFalse(),
			FileSystemPathEndParts:         test.EmptyFileSystemPathEndParts,
			DuplicateFilePathEndPartGroups: emptyPathEndPartGroups,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.Metadata.Name, func(t *testing.T) {
			// arrange and teardown
			directory := test.TestingCreateTempFileSystemStructureOrGetEmptyString(t, tc.FileSystemPathEndParts)
			defer test.TestingRemoveDirectoryTree(t, directory)
			var builder strings.Builder
			if len(tc.DuplicateFilePathEndPartGroups) > 0 {
				if len(tc.DuplicateFilePathEndPartGroups[0][0]) > 0 {
					duplicateFilePath := filepath.Join(directory, tc.DuplicateFilePathEndPartGroups[0][0])
					test.TestingWriteFileContentWithContentAndIndex(t, duplicateFilePath, 0)
					test.TestingWriteString(t, duplicateFilePath, &builder)
				}
				for i := 1; i < len(tc.DuplicateFilePathEndPartGroups[0]); i++ {
					testingWriteNewlineString(t, &builder)
					duplicateFilePath := filepath.Join(directory, tc.DuplicateFilePathEndPartGroups[0][i])
					test.TestingWriteFileContentWithContentAndIndex(t, duplicateFilePath, 0)
					test.TestingWriteString(t, duplicateFilePath, &builder)
				}
				for i := 1; i < len(tc.DuplicateFilePathEndPartGroups); i++ {
					testingWriteNewlineString(t, &builder)
					for _, duplicateFilePathEndPart := range tc.DuplicateFilePathEndPartGroups[i] {
						testingWriteNewlineString(t, &builder)
						duplicateFilePath := filepath.Join(directory, duplicateFilePathEndPart)
						test.TestingWriteFileContentWithContentAndIndex(t, duplicateFilePath, i)
						test.TestingWriteString(t, duplicateFilePath, &builder)
					}
				}
			}
			var fileSystemNodes []utils.FileSystemNode
			if directory != "" {
				fileSystemNodes = append(fileSystemNodes, utils.FileSystemNode{
					Path:        directory,
					IsDirectory: true,
				})
			}

			// act
			outcome, err := getDuplicateFilesAsNewlineSeparatedString(fileSystemNodes)

			// assert
			test.TestingAssertErrorToWantErrorAndOutcomeToBuilderString(t, err, tc.Metadata.WantErr, builder, outcome)
		})
	}
}
