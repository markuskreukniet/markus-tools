package test

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/markuskreukniet/markus-tools/go/utils"
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

// TODO: rename and use in TestPlainTextFilesToText?
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

func TestingCreateTempFileSystemStructureOrGetEmptyString(t *testing.T, fileSystemPathEndParts FileSystemPathEndParts) string {
	t.Helper()
	if len(fileSystemPathEndParts.DirectoryPathEndParts) == 0 {
		return ""
	}

	// Create a temporary file system structure.
	tempDirectory, err := os.MkdirTemp("", "markus-tools go test")
	if err != nil {
		t.Errorf("Failed to create the temporary directory: %v", err)
	}
	for _, part := range fileSystemPathEndParts.DirectoryPathEndParts {
		if err := os.MkdirAll(filepath.Join(tempDirectory, part), 0755); err != nil {
			t.Errorf("Failed to create directory in temporary directory: %v", err)
		}
	}
	for _, part := range fileSystemPathEndParts.FilePathEndParts {
		if err := os.WriteFile(filepath.Join(tempDirectory, part), []byte{}, 0666); err != nil {
			t.Errorf("Failed to create a file: %v", err)
		}
	}
	return tempDirectory
}

// TODO: comment
// type plainTextFile struct {
// 	name    string
// 	content string
// }
// type directoryWithOptionalFile struct {
// 	path          string
// 	timeModified  time.Time
// 	plainTextFile *plainTextFile
// }

func TestingTrimSpaceAndSplitOnComma(delimitedCommaString string) []string {
	return strings.Split(strings.TrimSpace(delimitedCommaString), ",")
}

func TestingCreateFilesAndDirectories(t *testing.T, input string) (string, []utils.FileSystemNode) {
	t.Helper()

	// if empty input string, return empty temporary directory file path and empty FileSystemNode slice
	if input == "" {
		return "", nil
	}

	// create temporary directory tree
	// return temporary directory file path and FileSystemNode slice
	tempDirectory, err := os.MkdirTemp("", "markus-tools go test")
	if err != nil {
		t.Errorf("Failed to create a temporary directory: %v", err)
	}
	var fileSystemNodes []utils.FileSystemNode
	for _, delimitedCommaString := range strings.Split(strings.TrimSuffix(strings.TrimSpace(input), ";"), ";") {
		directoryWithOptionalFileAsStrings := TestingTrimSpaceAndSplitOnComma(delimitedCommaString)
		filePath := filepath.Join(tempDirectory, filepath.FromSlash(directoryWithOptionalFileAsStrings[0]))
		exists, err := utils.FileOrDirectoryExists(filePath)
		if err != nil {
			t.Errorf("Failed to check if a file or directory exists: %v", err)
		}
		if !exists {
			if err := os.MkdirAll(filePath, 0755); err != nil {
				t.Errorf("Failed to create a directory in the temporary directory: %v", err)
			}
		}
		isDirectory := directoryWithOptionalFileAsStrings[2] == ""
		if !isDirectory {
			filePath = filepath.Join(filePath, directoryWithOptionalFileAsStrings[2])
			if err := os.WriteFile(filePath, []byte(directoryWithOptionalFileAsStrings[3]), 0666); err != nil {
				t.Errorf("Failed to create a file: %v", err)
			}
		}
		fileSystemNodes = append(fileSystemNodes, utils.FileSystemNode{
			Path:        filePath,
			IsDirectory: isDirectory,
		})
	}
	return tempDirectory, fileSystemNodes
}
