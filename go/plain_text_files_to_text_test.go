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
		Metadata                  test.TestCaseMetadata
		PlainTextFilePathEndParts []string
		FileSystemNodes           []utils.FileSystemNode
	}{
		{
			Metadata:                  test.TestingCreateTestCaseMetadataWithNameBasicAndWantErrFalse(),
			PlainTextFilePathEndParts: plainTextFilePathEndParts,
			FileSystemNodes:           fileSystemNodes,
		},
		{
			Metadata:                  test.TestingCreateTestCaseMetadataWithNameEmptyFileSystemNodesAndWantErrFalse(),
			PlainTextFilePathEndParts: test.EmptyPathEndParts,
			FileSystemNodes:           emptyFileSystemNodes,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.Metadata.Name, func(t *testing.T) {
			// arrange and teardown
			directory := test.TestingCreateTempFileSystemStructureOrGetEmptyString(t, fileSystemPathEndParts)
			defer test.TestingRemoveDirectoryTree(t, directory)
			for i := range tc.FileSystemNodes {
				tc.FileSystemNodes[i].Path = filepath.Join(directory, tc.FileSystemNodes[i].Path)
			}
			var builder strings.Builder
			// TODO: tc.PlainTextFilePathEndParts to full path happens also in testingCreateTempFileSystemStructureOrGetEmptyString?
			if len(tc.PlainTextFilePathEndParts) > 0 {
				fullPath := filepath.Join(directory, tc.PlainTextFilePathEndParts[0])
				content := testingCreateContentString(fullPath, 0)
				test.TestingWriteFileContent(t, fullPath, content)
				test.TestingWriteString(t, testingLastPathElementOnNewline(tc.PlainTextFilePathEndParts[0]), &builder)
				test.TestingWriteString(t, content, &builder)
				for i := 1; i < len(tc.PlainTextFilePathEndParts); i++ {
					fullPath := filepath.Join(directory, tc.PlainTextFilePathEndParts[i])
					content := testingCreateContentString(fullPath, i)
					test.TestingWriteFileContent(t, fullPath, content)
					test.TestingWriteString(t, "\n\n", &builder)
					test.TestingWriteString(t, testingLastPathElementOnNewline(tc.PlainTextFilePathEndParts[i]), &builder)
					test.TestingWriteString(t, content, &builder)
				}
			}

			// act
			outcome, err := plainTextFilesToText(tc.FileSystemNodes)

			// assert
			test.TestingAssertErrorToWantErrorAndOutcomeToBuilderString(t, err, tc.Metadata.WantErr, builder, outcome)
		})
	}
}
