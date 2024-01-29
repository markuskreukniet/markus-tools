package test

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"
)

type TestCaseMetadata struct {
	Name    string
	WantErr bool
}

type FileSystemPathEndParts struct {
	DirectoryPathEndParts []string
	FilePathEndParts      []string
}

// TODO: there are no files in the root (temp dir)
var (
	DirectoryEmpty               = "directory empty"
	Directory1                   = "directory 1"
	Directory2                   = "directory 2"
	Directory2WithDirectoryEmpty = filepath.Join(Directory2, DirectoryEmpty)
	Directory2WithDirectory3     = filepath.Join(Directory2, "directory 3")
	Directory2WithDirectory4     = filepath.Join(Directory2, "directory 4")

	TxtFile1 = filepath.Join(Directory1, "file 1.txt")
	TxtFile2 = filepath.Join(Directory1, "file 2.txt")
	TxtFile3 = filepath.Join(Directory2WithDirectory3, "file 3.txt")
	TxtFile4 = filepath.Join(Directory2WithDirectory3, "file 4.txt")
	TxtFile5 = filepath.Join(Directory2WithDirectory3, "file 5.txt")
	TxtFile6 = filepath.Join(Directory2WithDirectory4, "file 6.txt")

	JpgFile4 = filepath.Join(Directory1, "file 4.jpg")

	TxtFileNonExistent1 = "non existent 1.txt"

	EmptyPathEndParts []string

	EmptyFileSystemPathEndParts = FileSystemPathEndParts{
		DirectoryPathEndParts: EmptyPathEndParts,
		FilePathEndParts:      EmptyPathEndParts,
	}
)

func TestingCreateTestCaseMetadata(name string, wantErr bool) TestCaseMetadata {
	return TestCaseMetadata{
		Name:    name,
		WantErr: wantErr,
	}
}

func TestingCreateTestCaseMetadataWithWantErrTrue(name string) TestCaseMetadata {
	return TestingCreateTestCaseMetadata(name, true)
}

func TestingCreateTestCaseMetadataWithNameBasicAndWantErrFalse() TestCaseMetadata {
	return TestingCreateTestCaseMetadata("Basic", false)
}

func TestingCreateTestCaseMetadataWithNameEmptyFileSystemNodesAndWantErrFalse() TestCaseMetadata {
	return TestingCreateTestCaseMetadata("Empty FileSystemNodes", false)
}

func TestingWriteFileContentWithContentAndIndex(t *testing.T, filePath string, index int) string {
	t.Helper()
	writtenContent := fmt.Sprintf("content %d", index)
	TestingWriteFileContent(t, filePath, writtenContent)
	return writtenContent
}

func TestingWriteFileContent(t *testing.T, filePath string, content string) {
	t.Helper()
	if err := os.WriteFile(filePath, []byte(content), 0666); err != nil {
		t.Errorf("Failed to write file content: %v", err)
	}
}
