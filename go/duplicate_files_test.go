package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func testingWriteFileTestContent(t *testing.T, duplicateFilePath string, index int) {
	if err := os.WriteFile(duplicateFilePath, []byte(fmt.Sprintf("content %d", index)), 0666); err != nil {
		t.Errorf("Failed to write file content: %v", err)
	}
}

func testingWriteNewlineString(t *testing.T, builder *strings.Builder) {
	// TODO: move writeNewlineString to a different file
	_, err := writeNewlineString(builder)
	if err != nil {
		t.Errorf("writeNewlineString failed: %v", err)
	}
}

func testingWriteString(t *testing.T, stringToWrite string, builder *strings.Builder) {
	_, err := builder.WriteString(stringToWrite)
	if err != nil {
		t.Errorf("Failed to write string: %v", err)
	}
}

func TestGetDuplicateFilesAsNewlineSeparatedString(t *testing.T) {
	// arrange
	directoryPathEndParts := []string{directoryEmpty, directory1, directory2WithDirectoryEmpty, directory2WithDirectory3, directory2WithDirectory4}
	filePathEndParts := []string{txtFile1, txtFile2, txtFile3, txtFile4, txtFile5, txtFile6}
	duplicateFilePathEndPartGroups := [][]string{{txtFile2, txtFile3}, {txtFile4, txtFile5, txtFile6}}

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
	}

	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			// arrange and tear down
			// TODO: move createTempFileSystemStructureOrGetEmptyString to a different file
			directory, err := createTempFileSystemStructureOrGetEmptyString(tc.DirectoryPathEndParts, tc.FilePathEndParts)
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
					testingWriteFileTestContent(t, duplicateFilePath, 0)
					testingWriteString(t, duplicateFilePath, &builder)
				}
				for i := 1; i < len(tc.DuplicateFilePathEndPartGroups[0]); i++ {
					testingWriteNewlineString(t, &builder)
					duplicateFilePath := filepath.Join(directory, tc.DuplicateFilePathEndPartGroups[0][i])
					testingWriteFileTestContent(t, duplicateFilePath, 0)
					testingWriteString(t, duplicateFilePath, &builder)
				}
				for i := 1; i < len(tc.DuplicateFilePathEndPartGroups); i++ {
					testingWriteNewlineString(t, &builder)
					for _, duplicateFilePathEndPart := range tc.DuplicateFilePathEndPartGroups[i] {
						testingWriteNewlineString(t, &builder)
						duplicateFilePath := filepath.Join(directory, duplicateFilePathEndPart)
						testingWriteFileTestContent(t, duplicateFilePath, i)
						testingWriteString(t, duplicateFilePath, &builder)
					}
				}
			}
			fileSystemNodes := []FileSystemNode{{
				Path:        directory,
				IsDirectory: true,
			}}

			// act
			newlineSeparatedString, err := getDuplicateFilesAsNewlineSeparatedString(fileSystemNodes)

			// assert
			if (err != nil) != tc.WantErr {
				t.Fatalf("want error: %v, got %v", tc.WantErr, err)
			}
			if builder.String() != newlineSeparatedString {
				t.Fatalf("The newline-separated string is different than expected.")
			}
		})
	}
}
