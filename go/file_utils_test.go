package main

import (
	"os"
	"path/filepath"
	"testing"
)

func testingGetFileDetailFatalLogIfError(t *testing.T, err error) {
	// TODO: use t.Helper() also on other places?
	t.Helper()
	if err != nil {
		t.Fatalf("getFileDetail() error: %v", err)
	}
}

func testingFatalLogIfPathsAreNotEqual(t *testing.T, filePath string, fileDetailFilePath string) {
	t.Helper()
	if filePath != fileDetailFilePath {
		t.Errorf("Want Path %v, got %v", filePath, fileDetailFilePath)
	}
}

// TODO: change this test to a similar version as other tests
func TestGetFileDetail(t *testing.T) {
	// arrange
	directoryPathEndParts := []string{directory1}
	filePathEndParts := []string{txtFile1}

	// arrange and tear down
	directory := testingCreateTempFileSystemStructureOrGetEmptyString(t, directoryPathEndParts, filePathEndParts)
	defer func() {
		if err := os.RemoveAll(directory); err != nil {
			t.Errorf("Failed to remove the temporary directory: %v", err)
		}
	}()
	fullPath := filepath.Join(directory, filePathEndParts[0])
	writtenContent := testingWriteFileContentWithContentAndIndex(t, fullPath, 0)
	nonExistentFilePath := filepath.Join(directory, txtFileNonExistent1)

	// act
	dirDetail, err := getFileDetail(directory)
	testingGetFileDetailFatalLogIfError(t, err)
	fileDetail, err := getFileDetail(fullPath)
	testingGetFileDetailFatalLogIfError(t, err)
	_, err = getFileDetail(nonExistentFilePath)

	// assert
	// TODO: are all fileDetail properties checked?
	testingFatalLogIfPathsAreNotEqual(t, directory, dirDetail.Path)
	testingFatalLogIfPathsAreNotEqual(t, fullPath, fileDetail.Path)
	if fileDetail.Size != int64(len(writtenContent)) {
		t.Errorf("Want Size %v, got %v", len(writtenContent), fileDetail.Size)
	} else if err == nil {
		t.Errorf("Want an error when trying to get details of a non-existent file, but got none")
	}
}
