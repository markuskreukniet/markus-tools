package utils

import (
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"testing"
	"time"
)

type TestCaseMetadata struct {
	Name    string
	WantErr bool
}

type TestCaseInput struct {
	Metadata TestCaseMetadata
	Input    string
}

type InputLine []string

func (line InputLine) GetDirectoryPathPart() string {
	return line[0]
}

func (line InputLine) GetTimeModified() string {
	return line[1]
}

func (line InputLine) GetFileName() string {
	return line[2]
}

func (line InputLine) GetContent() string {
	return line[3]
}

func (line InputLine) HasNoContent() bool {
	return line.GetContent() != ""
}

func (line InputLine) JoinDirectoryPathPartWithFileName() string {
	return filepath.Join(line.GetDirectoryPathPart(), line.GetFileName())
}

func CreateInputLine(delimitedCommaString string) InputLine {
	return strings.Split(strings.TrimSpace(delimitedCommaString), ",")
}

func TestingCreateTestCaseMetadata(name string, wantErr bool) TestCaseMetadata {
	return TestCaseMetadata{
		Name:    name,
		WantErr: wantErr,
	}
}

func TestingCreateTestCaseInput(name, input string, wantErr bool) TestCaseInput {
	return TestCaseInput{
		Metadata: TestCaseMetadata{
			Name:    name,
			WantErr: wantErr,
		},
		Input: input,
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

// TODO: is it an arrange function? Should it be a separate function
func TestingTrimSpaceTrimSuffixSplitOnSemicolonAndSort(delimitedSemicolonString string) []string {
	delimitedCommaStrings := strings.Split(strings.TrimSuffix(strings.TrimSpace(delimitedSemicolonString), ";"), ";")
	slices.Sort(delimitedCommaStrings)
	return delimitedCommaStrings
}

// TODO: is it an arrange function? Should it be a separate function
// TODO: remove
func TestingTrimSpaceTrimSuffixOnSemicolonAndSplitOnSemicolon(delimitedSemicolonString string) []string {
	return strings.Split(strings.TrimSuffix(strings.TrimSpace(delimitedSemicolonString), ";"), ";")
}

// TODO: is it an arrange function?
func ToFilePathFromSlashAndJoin(filePath, filePathEndPart string) string {
	return filepath.Join(filePath, filepath.FromSlash(filePathEndPart))
}

func testingCreateDirectoryAll(t *testing.T, filePath string) {
	t.Helper()
	if err := os.MkdirAll(filePath, 0755); err != nil {
		t.Errorf("Failed to create a directory in the temporary directory: %v", err)
	}
}

func testingIfFileCreateFileAndAppendFileSystemNode(t *testing.T, isDirectory bool, filePath string, inputLine []string, fileSystemNodes *[]FileSystemNode) {
	t.Helper()
	if !isDirectory {
		filePath = filepath.Join(filePath, inputLine[2])
		if err := os.WriteFile(filePath, []byte(inputLine[3]), 0666); err != nil {
			t.Errorf("Failed to create a file: %v", err)
		}
		if inputLine[1] != "" {
			// 2006-01-02T15:04:05Z is ISO 8601 format
			timeModified, err := time.Parse("2006-01-02T15:04:05Z", inputLine[1])
			if err != nil {
				t.Errorf("Failed to parse time: %v", err)
			}
			if os.Chtimes(filePath, time.Now(), timeModified); err != nil {
				t.Errorf("Failed to change the access and modification times of the file: %v", err)
			}
		}
	}
	*fileSystemNodes = append(*fileSystemNodes, FileSystemNode{
		Path:        filePath,
		IsDirectory: isDirectory,
	})
}

// if empty input string, return empty temporary directory file path and empty FileSystemNode slice
func isInputEmpty(input string) bool {
	return input == ""
}

func createTemporaryDirectory(t *testing.T) string {
	t.Helper()
	temporaryDirectory, err := os.MkdirTemp("", "markus-tools go test")
	if err != nil {
		t.Errorf("Failed to create a temporary directory: %v", err)
	}
	return temporaryDirectory
}

// TODO: maybe using an [][][]string is not needed
// It should not always have to return a slice, but it is fine for testing.
// And disk I/O operations are significantly slower than in-memory operations.
func TestingCreateFilesAndDirectoriesByMultipleInputs(t *testing.T, input string) ([]string, []FileSystemNode) {
	t.Helper()
	if isInputEmpty(input) {
		return nil, nil
	}

	// create input groups
	var inputGroups [][][]string
	delimitedCommaStrings := TestingTrimSpaceTrimSuffixSplitOnSemicolonAndSort(input)
	for i := range delimitedCommaStrings {
		index := len(inputGroups) - 1
		inputLine := CreateInputLine(delimitedCommaStrings[i])

		// probably not optimal but results in less code, which is fine for testing
		if i == 0 || inputLine[0] == "" || inputLine[0] != inputGroups[index][0][0] {
			inputGroups = append(inputGroups, [][]string{inputLine})
		} else {
			inputGroups[index] = append(inputGroups[index], inputLine)
		}
	}

	// create and return temporary directories
	// create and return fileSystemNodes
	var tempDirectories []string
	var fileSystemNodes []FileSystemNode
	for _, group := range inputGroups {
		temporaryDirectory := createTemporaryDirectory(t)
		tempDirectories = append(tempDirectories, temporaryDirectory)
		for i, inputLine := range group {
			filePath := ToFilePathFromSlashAndJoin(temporaryDirectory, inputLine[0])
			isDirectory := inputLine[2] == ""

			// probably not optimal but results in less code, which is fine for testing
			// It should be possible to add more than one empty directory.
			if i == 0 || isDirectory {
				testingCreateDirectoryAll(t, filePath)
			}
			testingIfFileCreateFileAndAppendFileSystemNode(t, inputLine[2] == "", filePath, inputLine, &fileSystemNodes)
		}
	}
	return tempDirectories, fileSystemNodes
}

// This function has to stay for synchronizing directory trees.
// When we add a prefix to all input lines so that TestingCreateFilesAndDirectoriesByMultipleInputs can be used, all the folders with that prefix are added to the destination directory when syncing.
// It should not always have to return a slice, but it is fine for testing.
// And disk I/O operations are significantly slower than in-memory operations.
func TestingCreateFilesAndDirectoriesByOneInput(t *testing.T, input string) (string, []FileSystemNode) {
	t.Helper()
	if isInputEmpty(input) {
		return "", nil
	}

	temporaryDirectory := createTemporaryDirectory(t)
	var fileSystemNodes []FileSystemNode
	previousDirectoryFilePathPart := ""
	for _, delimitedCommaString := range TestingTrimSpaceTrimSuffixSplitOnSemicolonAndSort(input) {
		inputLine := CreateInputLine(delimitedCommaString)
		filePath := ToFilePathFromSlashAndJoin(temporaryDirectory, inputLine[0])
		if inputLine[0] != previousDirectoryFilePathPart {
			testingCreateDirectoryAll(t, filePath)

			// probably not optimal but results in less code, which is fine for testing
			previousDirectoryFilePathPart = inputLine[0]
		}
		testingIfFileCreateFileAndAppendFileSystemNode(t, inputLine[2] == "", filePath, inputLine, &fileSystemNodes)
	}
	return temporaryDirectory, fileSystemNodes
}
