package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/markuskreukniet/markus-tools/go/utils"
	"github.com/markuskreukniet/markus-tools/go/utils/test"
)

type plainTextFile struct {
	name    string
	content string
}

type directoryWithOptionalFile struct {
	path          string
	timeModified  time.Time
	plainTextFile *plainTextFile
}

func testingCreateFilesAndDirectories(t *testing.T, directoriesWithOptionalFileAsDelimitedSemicolonString string) string {
	t.Helper()

	//
	var directoriesWithOptionalFile []directoryWithOptionalFile
	directoriesWithOptionalFileAsDelimitedSemicolonString = strings.TrimSuffix(directoriesWithOptionalFileAsDelimitedSemicolonString, ";")
	directoriesWithOptionalFileAsDelimitedCommaString := strings.Split(directoriesWithOptionalFileAsDelimitedSemicolonString, ";")
	for _, delimitedCommaString := range directoriesWithOptionalFileAsDelimitedCommaString {
		directoryWithOptionalFileAsStrings := strings.Split(strings.TrimSpace(delimitedCommaString), ",")
		var file *plainTextFile = nil
		if directoryWithOptionalFileAsStrings[2] != "" {
			file = &plainTextFile{
				name:    directoryWithOptionalFileAsStrings[2],
				content: directoryWithOptionalFileAsStrings[3],
			}
		}
		// TODO: timeModified from string
		directoriesWithOptionalFile = append(directoriesWithOptionalFile, directoryWithOptionalFile{
			path:          directoryWithOptionalFileAsStrings[0],
			timeModified:  time.Now(),
			plainTextFile: file,
		})
	}

	//
	if len(directoriesWithOptionalFile) == 0 {
		return ""
	}
	tempDirectory, err := os.MkdirTemp("", "markus-tools go test")
	if err != nil {
		t.Errorf("Failed to create a temporary directory: %v", err)
	}
	for _, directoryWithOptionalFile := range directoriesWithOptionalFile {
		filePath := ""
		if directoryWithOptionalFile.path != "" {
			filePath = filepath.Join(tempDirectory, filepath.FromSlash(directoryWithOptionalFile.path))
			exists, err := utils.FileOrDirectoryExists(filePath)
			if err != nil {
				t.Errorf("Failed to check if a file or directory exists: %v", err)
			}
			if !exists {
				if err := os.MkdirAll(filePath, 0755); err != nil {
					t.Errorf("Failed to create a directory in the temporary directory: %v", err)
				}
			}
		}
		if directoryWithOptionalFile.plainTextFile != nil {
			if err := os.WriteFile(filepath.Join(filePath, directoryWithOptionalFile.plainTextFile.name),
				[]byte(directoryWithOptionalFile.plainTextFile.content), 0666); err != nil {
				t.Errorf("Failed to create a file: %v", err)
			}
		}
	}

	return tempDirectory
}

func TestFilesToDateRangeDirectory(t *testing.T) {
	// arrange
	// inputAsDelimitedString := `
	// 	,,txt 0.txt,;
	// 	empty,,nil,;
	// 	directory 1,,txt 1.txt,;
	// 	directory 1,,jpg 1.jpg,;
	// 	directory 2/directory 3,,txt 2 3.txt,;
	// `
	// destinationInputAsDelimitedString := `
	// 	,,txt 0.txt,;
	// 	2020-01-20,,nil,;
	// 	2020-02-20,,txt 02.txt,;
	// 	2020-03-20,,nil,;
	// 	2020-04-20 - 2020-04-21,,nil,;
	// 	2020-05-20 - 2020-05-21,,txt 05.txt,;
	// 	2020-06-20 - 2020-06-21,,nil,;
	// `
	// inputDirectoriesWithOptionalFile := createDirectoriesWithOptionalFile(inputAsDelimitedString)
	// destinationInputDirectoriesWithOptionalFile := createDirectoriesWithOptionalFile(destinationInputAsDelimitedString)

	// TODO: duplicate naming of files

	testCases := []struct {
		metadata test.TestCaseMetadata
	}{
		{
			metadata: test.TestingCreateTestCaseMetadataWithNameBasicAndWantErrFalse(),
		},
		{
			metadata: test.TestingCreateTestCaseMetadataWithNameEmptyFileSystemNodesAndWantErrFalse(),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.metadata.Name, func(t *testing.T) {
			// arrange and teardown
			// directoryTree := test.TestingCreateTempFileSystemStructureOrGetEmptyString(t, fileSystemPathEndParts)
			// defer test.TestingRemoveDirectoryTree(t, directoryTree)

			// act

			// assert
		})
	}
}
