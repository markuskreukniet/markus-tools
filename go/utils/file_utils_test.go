package utils

import (
	"path/filepath"
	"testing"
	"time"

	"github.com/markuskreukniet/markus-tools/go/utils/test"
)

func testingGetFileDetailLogIfError(t *testing.T, err error) {
	t.Helper()
	if err != nil {
		t.Errorf("getFileDetail() error: %v", err)
	}
}

// TODO: change this test to a similar version as other tests
func TestGetFileDetail(t *testing.T) {
	// arrange
	fileSystemPathEndParts := test.FileSystemPathEndParts{
		DirectoryPathEndParts: []string{test.Directory1},
		FilePathEndParts:      []string{test.TxtFile1},
	}

	// arrange and teardown
	directory := test.TestingCreateTempFileSystemStructureOrGetEmptyString(t, fileSystemPathEndParts)
	defer test.TestingRemoveDirectoryTree(t, directory)
	fullPath := filepath.Join(directory, fileSystemPathEndParts.FilePathEndParts[0])
	writtenContent := test.TestingWriteFileContentWithContentAndIndex(t, fullPath, 0)
	nonExistentFilePath := filepath.Join(directory, test.TxtFileNonExistent1)

	// act
	dirDetail, err := GetFileDetail(directory)
	testingGetFileDetailLogIfError(t, err)
	fileDetail, err := GetFileDetail(fullPath)
	testingGetFileDetailLogIfError(t, err)
	_, err = GetFileDetail(nonExistentFilePath)

	// assert
	test.TestingAssertEqualStrings(t, directory, dirDetail.Path)
	test.TestingAssertEqualStrings(t, fullPath, fileDetail.Path)
	if fileDetail.Size != int64(len(writtenContent)) {
		t.Errorf("Want Size %v, got %v", len(writtenContent), fileDetail.Size)
	}
	if err == nil {
		t.Errorf("Want an error when trying to get details of a non-existent file, but got none")
	}

	// Check if the file modification time is within the last minute, which is not optimal.
	currentTime := time.Now()
	if fileDetail.ModificationTime.Before(currentTime.Add(-time.Minute)) || fileDetail.ModificationTime.After(currentTime) {
		t.Errorf("Modification time %v is not within the expected range.", fileDetail.ModificationTime)
	}
}

// TODO: has copied and duplicate things and test cases are dirty
func TestFileOrDirectoryExists(t *testing.T) {
	// arrange
	fileSystemPathEndParts := test.FileSystemPathEndParts{
		DirectoryPathEndParts: []string{test.Directory1},
		FilePathEndParts:      []string{test.TxtFile1},
	}

	// arrange and teardown
	directory := test.TestingCreateTempFileSystemStructureOrGetEmptyString(t, fileSystemPathEndParts)
	defer test.TestingRemoveDirectoryTree(t, directory)

	// act
	dirExists, err := FileOrDirectoryExists(filepath.Join(directory, test.Directory1))
	if err != nil {
		t.Errorf("dirExists error")
	}
	fileExists, err := FileOrDirectoryExists(filepath.Join(directory, test.TxtFile1))
	if err != nil {
		t.Errorf("fileExists error")
	}
	fileDoesNotExists, err := FileOrDirectoryExists(filepath.Join(directory, test.TxtFileNonExistent1))
	if err != nil {
		t.Errorf("fileDoesNotExists error")
	}

	// assert
	if !dirExists || !fileExists || fileDoesNotExists {
		t.Errorf("TestFileOrDirectoryExists assert error")
	}
}
