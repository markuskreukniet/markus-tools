package main

import (
	"path/filepath"
	"testing"

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
	fileSystemPathEndParts := FileSystemPathEndParts{
		DirectoryPathEndParts: []string{directory1},
		FilePathEndParts:      []string{txtFile1},
	}

	// arrange and teardown
	directory := testingCreateTempFileSystemStructureOrGetEmptyString(t, fileSystemPathEndParts)
	defer test.TestRemoveDirectoryTree(t, directory)
	fullPath := filepath.Join(directory, fileSystemPathEndParts.FilePathEndParts[0])
	writtenContent := testingWriteFileContentWithContentAndIndex(t, fullPath, 0)
	nonExistentFilePath := filepath.Join(directory, txtFileNonExistent1)

	// act
	dirDetail, err := getFileDetail(directory)
	testingGetFileDetailLogIfError(t, err)
	fileDetail, err := getFileDetail(fullPath)
	testingGetFileDetailLogIfError(t, err)
	_, err = getFileDetail(nonExistentFilePath)

	// assert
	// TODO: are all fileDetail properties checked?
	test.TestingAssertEqualStrings(t, directory, dirDetail.Path)
	test.TestingAssertEqualStrings(t, fullPath, fileDetail.Path)
	if fileDetail.Size != int64(len(writtenContent)) {
		t.Errorf("Want Size %v, got %v", len(writtenContent), fileDetail.Size)
	} else if err == nil {
		t.Errorf("Want an error when trying to get details of a non-existent file, but got none")
	}
}
