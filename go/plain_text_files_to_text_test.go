package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func testingLastPathElementOnNewline(filePath string) string {
	return fmt.Sprintf("%s\n", filepath.Base(filePath))
}

func testingCreateContentString(filePath string, index int) string {
	return fmt.Sprintf("content %s %d 1\ncontent %s %d 2", filePath, index, filePath, index)
}

// TODO: there are duplicate or useless things, such as statements, strings, and structs, probably also in other tests
// TODO: PlainTextFilePathEndParts is useless, should use only FileSystemNodes?
func TestPlainTextFilesToText(t *testing.T) {
	// arrange
	directoryPathEndParts := []string{directory1, directory2WithDirectory3, directory2WithDirectory4}
	filePathEndParts := []string{txtFile1, txtFile3, txtFile6, jpgFile4}
	plainTextFilePathEndParts := []string{txtFile1, txtFile3, txtFile6}
	fileSystemNodes := []FileSystemNode{
		{
			Path:        txtFile1,
			IsDirectory: false,
		},
		{
			Path:        directory2,
			IsDirectory: true,
		},
	}
	var emptyFileSystemNodes []FileSystemNode

	testCases := []struct {
		Metadata                  TestCaseMetadata
		PlainTextFilePathEndParts []string
		FileSystemNodes           []FileSystemNode
	}{
		{
			Metadata:                  testingCreateTestCaseMetadataWithNameBasicAndWantErrFalse(),
			PlainTextFilePathEndParts: plainTextFilePathEndParts,
			FileSystemNodes:           fileSystemNodes,
		},
		{
			Metadata:                  testingCreateTestCaseMetadataWithNameEmptyFileSystemNodesAndWantErrFalse(),
			PlainTextFilePathEndParts: emptyPathEndParts,
			FileSystemNodes:           emptyFileSystemNodes,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.Metadata.Name, func(t *testing.T) {
			// arrange and tear down
			directory, err := testingCreateTempFileSystemStructureOrGetEmptyString(directoryPathEndParts, filePathEndParts)
			if err != nil {
				t.Fatalf("Failed to create the temporary directory: %v", err)
			}
			defer func() {
				if err := os.RemoveAll(directory); err != nil {
					t.Errorf("Failed to remove the temporary directory: %v", err)
				}
			}()
			for i := range tc.FileSystemNodes {
				tc.FileSystemNodes[i].Path = filepath.Join(directory, tc.FileSystemNodes[i].Path)
			}
			var builder strings.Builder
			// TODO: tc.PlainTextFilePathEndParts to full path happens also in testingCreateTempFileSystemStructureOrGetEmptyString?
			if len(tc.PlainTextFilePathEndParts) > 0 {
				fullPath := filepath.Join(directory, tc.PlainTextFilePathEndParts[0])
				content := testingCreateContentString(fullPath, 0)
				testingWriteFileContent(t, fullPath, content)
				testingWriteString(t, testingLastPathElementOnNewline(tc.PlainTextFilePathEndParts[0]), &builder)
				testingWriteString(t, content, &builder)
				for i := 1; i < len(tc.PlainTextFilePathEndParts); i++ {
					fullPath := filepath.Join(directory, tc.PlainTextFilePathEndParts[i])
					content := testingCreateContentString(fullPath, i)
					testingWriteFileContent(t, fullPath, content)
					testingWriteString(t, "\n\n", &builder)
					testingWriteString(t, testingLastPathElementOnNewline(tc.PlainTextFilePathEndParts[i]), &builder)
					testingWriteString(t, content, &builder)
				}
			}

			// act
			outcome, err := plainTextFilesToText(tc.FileSystemNodes)

			// assert
			testingAssertErrorToWantErrorAndOutcomeToBuilderString(t, err, tc.Metadata.WantErr, outcome, builder)
		})
	}
}
