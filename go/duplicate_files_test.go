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
	if _, err := utils.WriteNewlineString(builder); err != nil {
		t.Errorf("writeNewlineString error: %v", err)
	}
}

func TestGetDuplicateFilesAsNewlineSeparatedString(t *testing.T) {
	// content := "content 1\ncontent 1"
	// content2 := "content 2\ncontent 2"

	// input := `
	// 	empty,,,;
	// 	directory 1,,txt 1.txt,;
	// 	directory 1,,txt 1 2.txt,` + content + `;
	// 	directory 2/empty,,,;
	// 	directory 2/directory 3,,txt 2-3.txt,` + content + `;
	// 	directory 2/directory 3,,txt 2-3 2.txt,` + content2 + `;
	// 	directory 2/directory 3,,txt 2-3 3.txt,` + content2 + `;
	// 	directory 2/directory 4,,txt 2-4.txt,` + content2 + `;
	// `

	// arrange
	fileSystemPathEndParts := test.FileSystemPathEndParts{
		DirectoryPathEndParts: []string{test.DirectoryEmpty, test.Directory1, test.Directory2WithDirectoryEmpty, test.Directory2WithDirectory3, test.Directory2WithDirectory4},
		FilePathEndParts:      []string{test.TxtFile1, test.TxtFile2, test.TxtFile3, test.TxtFile4, test.TxtFile5, test.TxtFile6},
	}
	duplicateFilePathEndPartGroups := [][]string{{test.TxtFile2, test.TxtFile3}, {test.TxtFile4, test.TxtFile5, test.TxtFile6}}
	var emptyPathEndPartGroups [][]string

	testCases := []struct {
		metadata                       test.TestCaseMetadata
		fileSystemPathEndParts         test.FileSystemPathEndParts
		duplicateFilePathEndPartGroups [][]string
	}{
		{
			metadata:                       test.TestingCreateTestCaseMetadataWithNameBasicAndWantErrFalse(),
			fileSystemPathEndParts:         fileSystemPathEndParts,
			duplicateFilePathEndPartGroups: duplicateFilePathEndPartGroups,
		},
		{
			metadata:                       test.TestingCreateTestCaseMetadataWithNameEmptyFileSystemNodesAndWantErrFalse(),
			fileSystemPathEndParts:         test.EmptyFileSystemPathEndParts,
			duplicateFilePathEndPartGroups: emptyPathEndPartGroups,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.metadata.Name, func(t *testing.T) {
			// arrange and teardown
			directory := test.TestingCreateTempFileSystemStructureOrGetEmptyString(t, tc.fileSystemPathEndParts)
			defer test.TestingRemoveDirectoryTree(t, directory)
			var builder strings.Builder
			if len(tc.duplicateFilePathEndPartGroups) > 0 {
				// probably not optimal but results in less code, which is fine for testing
				for i, group := range duplicateFilePathEndPartGroups {
					for j, part := range group {
						check := i > 0 && j == 0
						if check || j > 0 {
							testingWriteNewlineString(t, &builder)
							if check {
								testingWriteNewlineString(t, &builder)
							}
						}
						duplicateFilePath := filepath.Join(directory, part)
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
			test.TestingAssertErrorToWantErrorAndOutcomeToBuilderString(t, err, tc.metadata.WantErr, builder, outcome)
		})
	}
}
