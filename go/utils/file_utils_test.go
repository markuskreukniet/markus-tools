package utils

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
	// TODO: are all fileDetail properties checked?
	test.TestingAssertEqualStrings(t, directory, dirDetail.Path)
	test.TestingAssertEqualStrings(t, fullPath, fileDetail.Path)
	if fileDetail.Size != int64(len(writtenContent)) {
		t.Errorf("Want Size %v, got %v", len(writtenContent), fileDetail.Size)
	} else if err == nil {
		t.Errorf("Want an error when trying to get details of a non-existent file, but got none")
	}
}
