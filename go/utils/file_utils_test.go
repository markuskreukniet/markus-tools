package utils

import (
	"testing"
)

// func testingGetFileDetailLogIfError(t *testing.T, err error) {
// 	t.Helper()
// 	if err != nil {
// 		t.Fatalf("getFileDetail() error: %v", err)
// 	}
// }

// func TestGetFileDetail(t *testing.T) {
// 	// arrange
// 	fileSystemPathEndParts := FileSystemPathEndParts{
// 		DirectoryPathEndParts: []string{Directory1},
// 		FilePathEndParts:      []string{TxtFile1},
// 	}

// 	// arrange and tear down
// 	directory := TestingCreateTempFileSystemStructureOrGetEmptyString(t, fileSystemPathEndParts)
// 	defer TestingRemoveDirectoryTree(t, directory)
// 	fullPath := filepath.Join(directory, fileSystemPathEndParts.FilePathEndParts[0])
// 	writtenContent := WriteFileWithContentAndIndex(t, fullPath, 0)
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
// 		t.Fatalf("Want Size %v, got %v", len(writtenContent), fileDetail.Size)
// 	}
// 	if err == nil {
// 		t.Fatalf("Want an error when trying to get details of a non-existent file, but got none")
// 	}

// 	// Check if the file modification time is within the last minute, which is not optimal.
// 	currentTime := time.Now()
// 	if fileDetail.ModificationTime.Before(currentTime.Add(-time.Minute)) || fileDetail.ModificationTime.After(currentTime) {
// 		t.Fatalf("Modification time %v is not within the expected range.", fileDetail.ModificationTime)
// 	}
// }

func TestFileExists(t *testing.T) {
	// arrange
	input := `
		,,txt 0.txt,;
		empty,,,;
		directory 1/empty,,,;
		directory 1,,txt 1.txt,;
	`
	testCases := []TestCaseBasicWithWriteInput{
		CreateTestCaseBasicWithWriteInput(CreateTestCaseBasic("Basic", input, "", false), true),
		CreateTestCaseBasicWithWriteInput(CreateTestCaseBasic("Empty Input", "", "", false), false),
	}

	for _, tc := range testCases {
		t.Run(tc.TestCaseBasic.Name, func(t *testing.T) {
			// arrange and tear down
			directory := WriteFilesBySingleInput(t, tc.TestCaseBasic.Input)
			defer TMustRemoveAll(t, directory)

			files := createFilesData(t, directory, tc.TestCaseBasic.Input)

			for _, file := range files {
				// act
				exists, err := FileExists(file.CompleteFileInfo.AbsolutePath)

				// assert
				TMustAssertError(t, err, tc.TestCaseBasic.WantErr)
				TMustAssertEqualBools(t, tc.WriteInput, exists)
			}
		})
	}
}
