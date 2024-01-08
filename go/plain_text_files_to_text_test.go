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

	testCases := []struct {
		Name                      string
		DirectoryPathEndParts     []string
		FilePathEndParts          []string
		PlainTextFilePathEndParts []string
		WantErr                   bool
	}{
		{
			Name:                      "Basic",
			DirectoryPathEndParts:     directoryPathEndParts,
			FilePathEndParts:          filePathEndParts,
			PlainTextFilePathEndParts: plainTextFilePathEndParts,
			WantErr:                   false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			// arrange and tear down
			directory, err := testingCreateTempFileSystemStructureOrGetEmptyString(tc.DirectoryPathEndParts, tc.FilePathEndParts)
			if err != nil {
				t.Fatalf("Failed to create the temporary directory: %v", err)
			}
			defer func() {
				if err := os.RemoveAll(directory); err != nil {
					t.Errorf("Failed to remove the temporary directory: %v", err)
				}
			}()
			for i := range fileSystemNodes {
				fileSystemNodes[i].Path = filepath.Join(directory, fileSystemNodes[i].Path)
			}
			var builder strings.Builder
			// TODO: plainTextFilePathEndParts to full path happens also in testingCreateTempFileSystemStructureOrGetEmptyString?
			if len(plainTextFilePathEndParts) > 0 {
				fullPath := filepath.Join(directory, plainTextFilePathEndParts[0])
				content := testingCreateContentString(fullPath, 0)
				testingWriteFileContent(t, fullPath, content)
				testingWriteString(t, testingLastPathElementOnNewline(plainTextFilePathEndParts[0]), &builder)
				testingWriteString(t, content, &builder)
				for i := 1; i < len(plainTextFilePathEndParts); i++ {
					fullPath := filepath.Join(directory, plainTextFilePathEndParts[i])
					content := testingCreateContentString(fullPath, i)
					testingWriteFileContent(t, fullPath, content)
					testingWriteString(t, "\n\n", &builder)
					testingWriteString(t, testingLastPathElementOnNewline(plainTextFilePathEndParts[i]), &builder)
					testingWriteString(t, content, &builder)
				}
			}

			// act
			outcome, err := plainTextFilesToText(fileSystemNodes)

			// assert
			testingAssertErrorToWantErrorAndOutcomeToBuilderString(t, err, tc.WantErr, outcome, builder)
		})
	}
}
