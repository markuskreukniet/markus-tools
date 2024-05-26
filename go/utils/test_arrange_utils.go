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

// InputLine
type InputLine struct {
	elements                      []string
	directoryPathPartWithFileName string
}

func (line InputLine) GetDirectoryPathPart() string {
	return line.elements[0]
}

func (line InputLine) GetTimeModified() string {
	return line.elements[1]
}

func (line InputLine) GetFileName() string {
	return line.elements[2]
}

func (line InputLine) GetContent() string {
	return line.elements[3]
}

func (line InputLine) HasContent() bool {
	return line.GetContent() != ""
}

func (line InputLine) GetDirectoryPathPartWithFileName() string {
	if line.directoryPathPartWithFileName == "" {
		line.directoryPathPartWithFileName = filepath.Join(line.GetDirectoryPathPart(), line.GetFileName())
	}
	return line.directoryPathPartWithFileName
}

func CreateInputLine(delimitedCommaString string) InputLine {
	return InputLine{
		elements:                      strings.Split(strings.TrimSpace(delimitedCommaString), ","),
		directoryPathPartWithFileName: "",
	}
}

// RawInputLines
type RawInputLines []string

func CreateSortedRawInputLines(delimitedSemicolonString string) RawInputLines {
	rawInputLines := strings.Split(strings.TrimSuffix(strings.TrimSpace(delimitedSemicolonString), ";"), ";")
	slices.Sort(rawInputLines)
	return rawInputLines
}

type TestCaseMetadata struct {
	Name    string
	WantErr bool
}

type TestCaseInput struct {
	Metadata TestCaseMetadata
	Input    string
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

func CreateTemporaryDirectory(t *testing.T) string {
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
	rawInputLines := CreateSortedRawInputLines(input)
	for i := range rawInputLines {
		index := len(inputGroups) - 1
		inputLine := CreateInputLine(rawInputLines[i])

		// probably not optimal but results in less code, which is fine for testing
		part := inputLine.GetDirectoryPathPart()
		if i == 0 || part == "" || part != inputGroups[index][0][0] {
			inputGroups = append(inputGroups, [][]string{inputLine.elements})
		} else {
			inputGroups[index] = append(inputGroups[index], inputLine.elements)
		}
	}

	// create and return temporary directories
	// create and return fileSystemNodes
	var tempDirectories []string
	var fileSystemNodes []FileSystemNode
	for _, group := range inputGroups {
		temporaryDirectory := CreateTemporaryDirectory(t)
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

	temporaryDirectory := CreateTemporaryDirectory(t)
	var fileSystemNodes []FileSystemNode
	previousDirectoryFilePathPart := ""
	for _, rawInputLine := range CreateSortedRawInputLines(input) {
		inputLine := CreateInputLine(rawInputLine)
		filePath := ToFilePathFromSlashAndJoin(temporaryDirectory, inputLine.elements[0])
		if inputLine.elements[0] != previousDirectoryFilePathPart {
			testingCreateDirectoryAll(t, filePath)

			// probably not optimal but results in less code, which is fine for testing
			previousDirectoryFilePathPart = inputLine.elements[0]
		}
		testingIfFileCreateFileAndAppendFileSystemNode(t, inputLine.elements[2] == "", filePath, inputLine.elements, &fileSystemNodes)
	}
	return temporaryDirectory, fileSystemNodes
}
