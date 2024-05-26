package utils

import (
	"path/filepath"
	"testing"
)

// func testingGetFileDetailLogIfError(t *testing.T, err error) {
// 	t.Helper()
// 	if err != nil {
// 		t.Errorf("getFileDetail() error: %v", err)
// 	}
// }

// func TestGetFileDetail(t *testing.T) {
// 	// arrange
// 	fileSystemPathEndParts := FileSystemPathEndParts{
// 		DirectoryPathEndParts: []string{Directory1},
// 		FilePathEndParts:      []string{TxtFile1},
// 	}

// 	// arrange and teardown
// 	directory := TestingCreateTempFileSystemStructureOrGetEmptyString(t, fileSystemPathEndParts)
// 	defer TestingRemoveDirectoryTree(t, directory)
// 	fullPath := filepath.Join(directory, fileSystemPathEndParts.FilePathEndParts[0])
// 	writtenContent := TestingWriteFileContentWithContentAndIndex(t, fullPath, 0)
// 	nonExistentFilePath := filepath.Join(directory, TxtFileNonExistent1)

// 	// act
// 	dirDetail, err := GetFileDetail(directory)
// 	testingGetFileDetailLogIfError(t, err)
// 	fileDetail, err := GetFileDetail(fullPath)
// 	testingGetFileDetailLogIfError(t, err)
// 	_, err = GetFileDetail(nonExistentFilePath)

// 	// assert
// 	TestingAssertEqualStrings(t, directory, dirDetail.Path)
// 	TestingAssertEqualStrings(t, fullPath, fileDetail.Path)
// 	if fileDetail.Size != int64(len(writtenContent)) {
// 		t.Errorf("Want Size %v, got %v", len(writtenContent), fileDetail.Size)
// 	}
// 	if err == nil {
// 		t.Errorf("Want an error when trying to get details of a non-existent file, but got none")
// 	}

// 	// Check if the file modification time is within the last minute, which is not optimal.
// 	currentTime := time.Now()
// 	if fileDetail.ModificationTime.Before(currentTime.Add(-time.Minute)) || fileDetail.ModificationTime.After(currentTime) {
// 		t.Errorf("Modification time %v is not within the expected range.", fileDetail.ModificationTime)
// 	}
// }

func TestFileOrDirectoryExists(t *testing.T) {
	// arrange
	input := `
		,,txt 0.txt,;
		empty,,,;
		directory 1/empty,,,;
		directory 1,,txt 1.txt,;
	`
	testCases := []struct {
		testCaseInput    TestCaseInput
		inputToDirectory bool
	}{
		{
			testCaseInput:    TestingCreateTestCaseInput("Basic", input, false),
			inputToDirectory: true,
		},
		{
			testCaseInput:    TestingCreateTestCaseInput("Empty Input", "", false),
			inputToDirectory: false,
		},
	}

	// run testCases
	for _, tc := range testCases {
		t.Run(tc.testCaseInput.Metadata.Name, func(t *testing.T) {
			// arrange and teardown
			directory, _ := TestingCreateFilesAndDirectoriesByOneInput(t, tc.testCaseInput.Input)
			defer TestingRemoveDirectoryTree(t, directory)
			rawInputLines := CreateSortedRawInputLines(tc.testCaseInput.Input)
			for _, rawInputLine := range rawInputLines {
				filePath := directory
				if directory != "" {
					filePath = filepath.Join(directory, CreateInputLine(rawInputLine).GetDirectoryPathPartWithFileName())
				}

				// act
				exists, err := FileOrDirectoryExists(filePath)
				if err != nil {
					t.Errorf("FileOrDirectoryExists error")
				}

				// assert
				if exists != tc.inputToDirectory {
					t.Errorf("A directory or file should exist, but it does not.")
				}
			}
		})
	}
}
