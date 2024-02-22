package main

import (
	"fmt"
	"path/filepath"
	"strings"
	"testing"

	"github.com/markuskreukniet/markus-tools/go/utils"
	"github.com/markuskreukniet/markus-tools/go/utils/test"
)

func testingLastPathElementOnNewline(filePath string) string {
	return fmt.Sprintf("%s\n", filepath.Base(filePath))
}

func testingCreateContentString(filePath string, index int) string {
	return fmt.Sprintf("content %s %d 1\ncontent %s %d 2", filePath, index, filePath, index)
}

func TestPlainTextFilesToText(t *testing.T) {
	// arrange
	fileSystemPathEndParts := test.FileSystemPathEndParts{
		DirectoryPathEndParts: []string{test.Directory1, test.Directory2WithDirectory3, test.Directory2WithDirectory4},
		FilePathEndParts:      []string{test.TxtFile1, test.TxtFile3, test.TxtFile6, test.JpgFile4},
	}
	plainTextFilePathEndParts := []string{test.TxtFile1, test.TxtFile3, test.TxtFile6}
	fileSystemNodes := []utils.FileSystemNode{
		{
			Path:        test.TxtFile1,
			IsDirectory: false,
		},
		{
			Path:        test.Directory2,
			IsDirectory: true,
		},
	}
	var emptyFileSystemNodes []utils.FileSystemNode

	testCases := []struct {
		metadata                  test.TestCaseMetadata
		plainTextFilePathEndParts []string
		fileSystemNodes           []utils.FileSystemNode
	}{
		{
			metadata:                  test.TestingCreateTestCaseMetadataWithNameBasicAndWantErrFalse(),
			plainTextFilePathEndParts: plainTextFilePathEndParts,
			fileSystemNodes:           fileSystemNodes,
		},
		{
			metadata:                  test.TestingCreateTestCaseMetadataWithNameEmptyFileSystemNodesAndWantErrFalse(),
			plainTextFilePathEndParts: test.EmptyPathEndParts,
			fileSystemNodes:           emptyFileSystemNodes,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.metadata.Name, func(t *testing.T) {
			// arrange and teardown
			directory := test.TestingCreateTempFileSystemStructureOrGetEmptyString(t, fileSystemPathEndParts)
			defer test.TestingRemoveDirectoryTree(t, directory)
			for i := range tc.fileSystemNodes {
				tc.fileSystemNodes[i].Path = filepath.Join(directory, tc.fileSystemNodes[i].Path)
			}
			var builder strings.Builder
			// TODO: tc.PlainTextFilePathEndParts to full path happens also in testingCreateTempFileSystemStructureOrGetEmptyString?
			if len(tc.plainTextFilePathEndParts) > 0 {
				fullPath := filepath.Join(directory, tc.plainTextFilePathEndParts[0])
				content := testingCreateContentString(fullPath, 0)
				test.TestingWriteFileContent(t, fullPath, content)
				test.TestingWriteString(t, testingLastPathElementOnNewline(tc.plainTextFilePathEndParts[0]), &builder)
				test.TestingWriteString(t, content, &builder)
				for i := 1; i < len(tc.plainTextFilePathEndParts); i++ {
					fullPath := filepath.Join(directory, tc.plainTextFilePathEndParts[i])
					content := testingCreateContentString(fullPath, i)
					test.TestingWriteFileContent(t, fullPath, content)
					test.TestingWriteString(t, "\n\n", &builder)
					test.TestingWriteString(t, testingLastPathElementOnNewline(tc.plainTextFilePathEndParts[i]), &builder)
					test.TestingWriteString(t, content, &builder)
				}
			}

			// act
			outcome, err := plainTextFilesToText(tc.fileSystemNodes)

			// assert
			test.TestingAssertErrorToWantErrorAndOutcomeToBuilderString(t, err, tc.metadata.WantErr, builder, outcome)
		})
	}
}
