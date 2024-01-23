package main

import (
	"path/filepath"
	"testing"
)

// TODO: rename Fatal also on other places
func testingGetFileDetailFatalLogIfError(t *testing.T, err error) {
	// TODO: use t.Helper() also on other places?
	t.Helper()
	if err != nil {
		t.Errorf("getFileDetail() error: %v", err)
	}
}

// TODO: change this test to a similar version as other tests
func TestGetFileDetail(t *testing.T) {
	// arrange
	directoryPathEndParts := []string{directory1}
	filePathEndParts := []string{txtFile1}

	// arrange and teardown
	directory := testingCreateTempFileSystemStructureOrGetEmptyString(t, directoryPathEndParts, filePathEndParts)
	defer testingRemoveDirectoryTree(t, directory)
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
	testingAssertEqualStrings(t, directory, dirDetail.Path)
	testingAssertEqualStrings(t, fullPath, fileDetail.Path)
	if fileDetail.Size != int64(len(writtenContent)) {
		t.Errorf("Want Size %v, got %v", len(writtenContent), fileDetail.Size)
	} else if err == nil {
		t.Errorf("Want an error when trying to get details of a non-existent file, but got none")
	}
}
