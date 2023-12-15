package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func writeFileContentAndWriteString(t *testing.T, directory string, filePathEndPart string, index int, builder *strings.Builder) {
	duplicateFilePath := filepath.Join(directory, filePathEndPart)
	if err := os.WriteFile(duplicateFilePath, []byte(fmt.Sprintf("content %d", index)), 0666); err != nil {
		t.Errorf("Failed to write file content: %v", err)
	}
	_, err := builder.WriteString(duplicateFilePath)
	if err != nil {
		t.Errorf("Failed to write string: %v", err)
	}
}

func TestGetDuplicateFilesAsNewlineSeparatedString(t *testing.T) {
	// arrange
	// TODO: copied and expanded file path structure from TestSynchronizeDirectoryTrees.
	// This structure should be part of an object which can return PathEndParts.
	directoryEmpty := "directory empty"
	directory1 := "directory 1"
	directory2 := "directory 2"
	directory2WithDirectoryEmpty := filepath.Join(directory2, directoryEmpty)
	directory2WithDirectory3 := filepath.Join(directory2, "directory 3")
	directory2WithDirectory4 := filepath.Join(directory2, "directory 4")

	txtFile1 := filepath.Join(directory1, "file 1.txt")
	txtFile2 := filepath.Join(directory1, "file 2.txt")
	txtFile3 := filepath.Join(directory2WithDirectory3, "file 3.txt")
	txtFile4 := filepath.Join(directory2WithDirectory3, "file 4.txt")
	txtFile5 := filepath.Join(directory2WithDirectory3, "file 5.txt")
	txtFile6 := filepath.Join(directory2WithDirectory4, "file 6.txt")

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
				for _, duplicateFilePathEndPart := range tc.DuplicateFilePathEndPartGroups[0] {
					writeFileContentAndWriteString(t, directory, duplicateFilePathEndPart, 0, &builder)
				}
				for i := 1; i < len(tc.DuplicateFilePathEndPartGroups); i++ {
					// TODO: move writeNewlineString to a different file
					_, err = writeNewlineString(&builder)
					if err != nil {
						t.Errorf("writeNewlineString failed: %v", err)
					}
					for _, duplicateFilePathEndPart := range tc.DuplicateFilePathEndPartGroups[i] {
						writeFileContentAndWriteString(t, directory, duplicateFilePathEndPart, i, &builder)
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
			// TODO: fmt.Printlns
			fmt.Println(builder.String())
			fmt.Println("--- --- ---")
			fmt.Println(newlineSeparatedString)
			if builder.String() != newlineSeparatedString {
				t.Fatalf("The newline-separated string is different than expected.")
			}
		})
	}
}
