package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func testingWriteFileContentWithContentAndIndex(t *testing.T, filePath string, index int) {
	testingWriteFileContent(t, filePath, fmt.Sprintf("content %d", index))
}

func TestGetDuplicateFilesAsNewlineSeparatedString(t *testing.T) {
	// arrange
	directoryPathEndParts := []string{directoryEmpty, directory1, directory2WithDirectoryEmpty, directory2WithDirectory3, directory2WithDirectory4}
	filePathEndParts := []string{txtFile1, txtFile2, txtFile3, txtFile4, txtFile5, txtFile6}
	duplicateFilePathEndPartGroups := [][]string{{txtFile2, txtFile3}, {txtFile4, txtFile5, txtFile6}}
	var emptyPathEndPartGroups [][]string

	testCases := []struct {
		Name                           string
		DirectoryPathEndParts          []string
		FilePathEndParts               []string
		DuplicateFilePathEndPartGroups [][]string
		WantErr                        bool
	}{
		{
			Name:                           "Basic",
			DirectoryPathEndParts:          directoryPathEndParts,
			FilePathEndParts:               filePathEndParts,
			DuplicateFilePathEndPartGroups: duplicateFilePathEndPartGroups,
			WantErr:                        false,
		},
		{
			Name:                           "Empty FileSystemNodes",
			DirectoryPathEndParts:          emptyPathEndParts,
			FilePathEndParts:               emptyPathEndParts,
			DuplicateFilePathEndPartGroups: emptyPathEndPartGroups,
			WantErr:                        false,
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
			testingAssertErrorToWantErrorAndOutcomeToBuilderString(t, err, tc.WantErr, outcome, builder)
		})
	}
}
