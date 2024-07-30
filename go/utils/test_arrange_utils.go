package utils

import (
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"sort"
	"strings"
	"testing"
	"time"
	"unicode"
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

func (line InputLine) IsDirectory() bool {
	return line.GetFileName() == ""
}

func (line InputLine) GetDirectoryPathPartWithFileName() string {
	if line.directoryPathPartWithFileName == "" {
		line.directoryPathPartWithFileName = filepath.Join(line.GetDirectoryPathPart(), line.GetFileName())
	}
	return line.directoryPathPartWithFileName
}

func CreateInputLine(delimitedCommaString string) InputLine {
	return InputLine{
		elements:                      strings.Split(delimitedCommaString, ","),
		directoryPathPartWithFileName: "",
	}
}

// RawInputLines
type RawInputLines []string

func createFileSystemFileByInputLine(t *testing.T, directoryPath, inputLine string) FileSystemFile {
	t.Helper()

	fields := strings.Split(inputLine, ",")

	directoryPath = ToFilePathFromSlashAndJoin(directoryPath, fields[0])
	data := fields[3]
	name := fields[2]
	filePath := filepath.Join(directoryPath, name)
	isDirectory := name == ""

	var timeModified time.Time
	if fields[1] != "" {
		timeModified = TestingParseTime(t, fields[1])
	}

	return CreateFileSystemFile(data, filePath, CreateFileMetadata(name, directoryPath, timeModified, 0, isDirectory))
}

func createSortedFileSystemFiles(t *testing.T, directoryPath, rawDelimitedSemicolonString string) []FileSystemFile {
	t.Helper()

	var files []FileSystemFile
	var inputLine []rune
	isCreatingInputLine := false

	rawDelimitedSemicolonString = strings.TrimSpace(rawDelimitedSemicolonString)

	for _, r := range rawDelimitedSemicolonString {
		if isCreatingInputLine {
			if r != ';' {
				inputLine = append(inputLine, r)
			} else {
				files = append(files, createFileSystemFileByInputLine(t, directoryPath, string(inputLine)))
				inputLine = nil
				isCreatingInputLine = false
			}
		} else if !unicode.IsSpace(r) {
			inputLine = append(inputLine, r)
			isCreatingInputLine = true
		}
	}

	sort.Slice(files, func(i, j int) bool {
		return files[i].Path < files[j].Path
	})

	return files
}

func CreateSortedRawInputLines(rawDelimitedSemicolonString string) []string {
	var inputLines []string
	var inputLine []rune
	isCreatingInputLine := false

	rawDelimitedSemicolonString = strings.TrimSpace(rawDelimitedSemicolonString)

	for _, r := range rawDelimitedSemicolonString {
		if isCreatingInputLine {
			if r != ';' {
				inputLine = append(inputLine, r)
			} else {
				inputLines = append(inputLines, string(inputLine))
				inputLine = nil
				isCreatingInputLine = false
			}
		} else if !unicode.IsSpace(r) {
			inputLine = append(inputLine, r)
			isCreatingInputLine = true
		}
	}

	slices.Sort(inputLines)

	return inputLines
}

// TestCase
type TestCaseMetadata struct {
	Name    string
	WantErr bool
}

type TestCaseInput struct {
	Metadata TestCaseMetadata
	Input    string
}

func CreateTestCaseMetadata(name string, wantErr bool) TestCaseMetadata {
	return TestCaseMetadata{
		Name:    name,
		WantErr: wantErr,
	}
}

func CreateTestCaseInput(name, input string, wantErr bool) TestCaseInput {
	return TestCaseInput{
		Metadata: TestCaseMetadata{
			Name:    name,
			WantErr: wantErr,
		},
		Input: input,
	}
}

func CreateTestCaseMetadataWithWantErrTrue(name string) TestCaseMetadata {
	return CreateTestCaseMetadata(name, true)
}

func CreateTestCaseMetadataWithNameBasicAndWantErrFalse() TestCaseMetadata {
	return CreateTestCaseMetadata("Basic", false)
}

// TODO: wrong naming
func TestingWriteFileWithContentAndIndex(t *testing.T, filePath string, index int) string {
	t.Helper()
	writtenContent := fmt.Sprintf("content %d", index)
	TestingWriteFile(t, filePath, writtenContent)
	return writtenContent
}

// TODO: should only receive t and file?
func TestingWriteFile(t *testing.T, filePath string, content string) {
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

// TODO: is it an arrange function?
// TODO: wrong naming, testing forgotten
func ToFilePathFromSlashAndJoin(filePath, filePathEndPart string) string {
	return filepath.Join(filePath, filepath.FromSlash(filePathEndPart))
}

func TestingCreateDirectoryAll(t *testing.T, filePath string) {
	t.Helper()
	if err := os.MkdirAll(filePath, 0755); err != nil {
		t.Errorf("Failed to create a directory in the temporary directory: %v", err)
	}
}

func TestingParseTime(t *testing.T, timeString string) time.Time {
	t.Helper()
	parsedTime, err := time.Parse(time.RFC3339, timeString)
	if err != nil {
		t.Errorf("Failed to parse time: %v", err)
	}
	return parsedTime
}

func testingIfFileWriteItAndAppendFileSystemNode(t *testing.T, file FileSystemFile, nodes *[]FileSystemNode) {
	t.Helper()

	if !file.FileMetadata.IsDirectory {
		TestingWriteFile(t, file.Path, file.Data)
		if !file.FileMetadata.TimeModified.IsZero() {
			if err := os.Chtimes(file.Path, time.Now(), file.FileMetadata.TimeModified); err != nil {
				t.Errorf("Failed to change the access and modification times of the file: %v", err)
			}
		}
	}

	*nodes = append(*nodes, FileSystemNode{
		Path:        file.Path,
		IsDirectory: file.FileMetadata.IsDirectory,
	})
}

func testingIfFileCreateFileAndAppendFileSystemNode(t *testing.T, filePath string, line InputLine, fileSystemNodes *[]FileSystemNode) {
	t.Helper()
	if !line.IsDirectory() {
		filePath = filepath.Join(filePath, line.GetFileName())
		TestingWriteFile(t, filePath, line.GetContent())
		if line.GetTimeModified() != "" {
			timeModified := TestingParseTime(t, line.GetTimeModified())
			if err := os.Chtimes(filePath, time.Now(), timeModified); err != nil {
				t.Errorf("Failed to change the access and modification times of the file: %v", err)
			}
		}
	}
	*fileSystemNodes = append(*fileSystemNodes, FileSystemNode{
		Path:        filePath,
		IsDirectory: line.IsDirectory(),
	})
}

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

func toRootDirectoryPath(filePath string) string {
	cleanPath := filepath.Clean(filePath)

	for {
		parentPath := filepath.Dir(cleanPath)
		if parentPath == "." || parentPath == cleanPath {
			break
		}
		cleanPath = parentPath
	}

	return cleanPath
}

func TestingWriteFilesByMultipleInputs(t *testing.T, input string) ([]string, []FileSystemNode) {
	t.Helper()

	if isInputEmpty(input) {
		return nil, nil
	}

	files := createSortedFileSystemFiles(t, "", input)
	length := len(files)

	if length == 0 {
		return nil, nil
	}

	fileGroups := [][]FileSystemFile{{files[0]}}
	previousRootDirectoryPath := toRootDirectoryPath(files[0].FileMetadata.DirectoryPath)
	index := 0

	for i := 1; i < length; i++ {
		rootDirectoryPath := toRootDirectoryPath(files[i].FileMetadata.DirectoryPath)
		if rootDirectoryPath == "." || rootDirectoryPath != previousRootDirectoryPath {
			fileGroups = append(fileGroups, []FileSystemFile{files[i]})
			previousRootDirectoryPath = rootDirectoryPath
			index++
		} else {
			fileGroups[index] = append(fileGroups[index], files[i])
		}
	}

	var temporaryDirectories []string
	var fileSystemNodes []FileSystemNode
	var previousDirectoryPath string

	for _, group := range fileGroups {
		directory := CreateTemporaryDirectory(t)
		temporaryDirectories = append(temporaryDirectories, directory)
		for i, file := range group {
			file.FileMetadata.DirectoryPath = filepath.Join(directory, file.FileMetadata.DirectoryPath)
			file.Path = filepath.Join(directory, file.Path)
			if i == 0 || file.FileMetadata.DirectoryPath != previousDirectoryPath {
				TestingCreateDirectoryAll(t, file.FileMetadata.DirectoryPath)
			}
			previousDirectoryPath = file.FileMetadata.DirectoryPath
			testingIfFileWriteItAndAppendFileSystemNode(t, file, &fileSystemNodes)
		}
	}

	return temporaryDirectories, fileSystemNodes
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
				TestingCreateDirectoryAll(t, filePath)
			}

			// TODO: should have been already an InputLine
			line := InputLine{
				elements:                      inputLine,
				directoryPathPartWithFileName: "",
			}

			testingIfFileCreateFileAndAppendFileSystemNode(t, filePath, line, &fileSystemNodes)
		}
	}
	return tempDirectories, fileSystemNodes
}

// When we add a prefix to all input lines so that TestingWriteFilesByMultipleInputs can be used, all the folders with that prefix are added to the destination directory when syncing.
// It should not always have to return a slice, but it is fine for testing.
// And disk I/O operations are significantly slower than in-memory operations.
func TestingWriteFilesByOneInput(t *testing.T, input string) (string, []FileSystemNode) {
	t.Helper()

	if isInputEmpty(input) {
		return "", nil
	}

	var nodes []FileSystemNode
	var previousDirectoryPath string
	directory := CreateTemporaryDirectory(t)
	files := createSortedFileSystemFiles(t, directory, input)

	for i := range files {
		if previousDirectoryPath != files[i].FileMetadata.DirectoryPath {
			TestingCreateDirectoryAll(t, files[i].FileMetadata.DirectoryPath)
			previousDirectoryPath = files[i].FileMetadata.DirectoryPath
		}
		testingIfFileWriteItAndAppendFileSystemNode(t, files[i], &nodes)
	}

	return directory, nodes
}
