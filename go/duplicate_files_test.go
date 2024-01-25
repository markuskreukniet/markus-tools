package main

import (
	"path/filepath"
	"strings"
	"testing"

	"github.com/markuskreukniet/markus-tools/go/utils/test"
)

func testingWriteNewlineString(t *testing.T, builder *strings.Builder) {
	t.Helper()
	_, err := writeNewlineString(builder)
	if err != nil {
		t.Errorf("writeNewlineString error: %v", err)
	}
}

func TestGetDuplicateFilesAsNewlineSeparatedString(t *testing.T) {
	// arrange
	fileSystemPathEndParts := FileSystemPathEndParts{
		DirectoryPathEndParts: []string{directoryEmpty, directory1, directory2WithDirectoryEmpty, directory2WithDirectory3, directory2WithDirectory4},
		FilePathEndParts:      []string{txtFile1, txtFile2, txtFile3, txtFile4, txtFile5, txtFile6},
	}
	duplicateFilePathEndPartGroups := [][]string{{txtFile2, txtFile3}, {txtFile4, txtFile5, txtFile6}}
	var emptyPathEndPartGroups [][]string

	testCases := []struct {
		Metadata                       TestCaseMetadata
		FileSystemPathEndParts         FileSystemPathEndParts
		DuplicateFilePathEndPartGroups [][]string
	}{
		{
			Metadata:                       testingCreateTestCaseMetadataWithNameBasicAndWantErrFalse(),
			FileSystemPathEndParts:         fileSystemPathEndParts,
			DuplicateFilePathEndPartGroups: duplicateFilePathEndPartGroups,
		},
		{
			Metadata:                       testingCreateTestCaseMetadataWithNameEmptyFileSystemNodesAndWantErrFalse(),
			FileSystemPathEndParts:         emptyFileSystemPathEndParts,
			DuplicateFilePathEndPartGroups: emptyPathEndPartGroups,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.Metadata.Name, func(t *testing.T) {
			// arrange and teardown
			directory := testingCreateTempFileSystemStructureOrGetEmptyString(t, tc.FileSystemPathEndParts)
			defer test.TestRemoveDirectoryTree(t, directory)
			var builder strings.Builder
			if len(tc.DuplicateFilePathEndPartGroups) > 0 {
				if len(tc.DuplicateFilePathEndPartGroups[0][0]) > 0 {
					duplicateFilePath := filepath.Join(directory, tc.DuplicateFilePathEndPartGroups[0][0])
					testingWriteFileContentWithContentAndIndex(t, duplicateFilePath, 0)
					testingWriteString(t, duplicateFilePath, &builder)
				}
				for i := 1; i < len(tc.DuplicateFilePathEndPartGroups[0]); i++ {
					testingWriteNewlineString(t, &builder)
					duplicateFilePath := filepath.Join(directory, tc.DuplicateFilePathEndPartGroups[0][i])
					testingWriteFileContentWithContentAndIndex(t, duplicateFilePath, 0)
					testingWriteString(t, duplicateFilePath, &builder)
				}
				for i := 1; i < len(tc.DuplicateFilePathEndPartGroups); i++ {
					testingWriteNewlineString(t, &builder)
					for _, duplicateFilePathEndPart := range tc.DuplicateFilePathEndPartGroups[i] {
						testingWriteNewlineString(t, &builder)
						duplicateFilePath := filepath.Join(directory, duplicateFilePathEndPart)
						testingWriteFileContentWithContentAndIndex(t, duplicateFilePath, i)
						testingWriteString(t, duplicateFilePath, &builder)
					}
				}
			}
			var fileSystemNodes []FileSystemNode
			if directory != "" {
				fileSystemNodes = append(fileSystemNodes, FileSystemNode{
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
