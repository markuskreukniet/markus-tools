package utils

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"testing"
	"time"
	"unicode"
)

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

	return CreateFileSystemFile(data, CreateFileMetadata(name, directoryPath, filePath, timeModified, 0, isDirectory))
}

func CreateSortedFileSystemFiles(t *testing.T, directoryPath, rawDelimitedSemicolonString string) []FileSystemFile {
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
		return files[i].FileMetadata.Path < files[j].FileMetadata.Path
	})

	return files
}

type TestCaseBasicDoubleInput struct {
	TestCaseBasic TestCaseBasic
	SecondInput   string
}

func CreateTestCaseBasicDoubleInput(testCaseBasic TestCaseBasic, secondInput string) TestCaseBasicDoubleInput {
	return TestCaseBasicDoubleInput{
		TestCaseBasic: testCaseBasic,
		SecondInput:   secondInput,
	}
}

type TestCaseBasicWithWriteInput struct {
	TestCaseBasic TestCaseBasic
	WriteInput    bool
}

func CreateTestCaseBasicWithWriteInput(testCaseBasic TestCaseBasic, writeInput bool) TestCaseBasicWithWriteInput {
	return TestCaseBasicWithWriteInput{
		TestCaseBasic: testCaseBasic,
		WriteInput:    writeInput,
	}
}

// TestCase
type TestCaseBasic struct {
	Name          string
	Input         string
	WantedOutcome string
	WantErr       bool
}

func CreateTestCaseBasic(name, input, wantedOutcome string, wantErr bool) TestCaseBasic {
	return TestCaseBasic{
		Name:          name,
		Input:         input,
		WantedOutcome: wantedOutcome,
		WantErr:       wantErr,
	}
}

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
		TestingWriteFile(t, file.FileMetadata.Path, file.Data)
		if !file.FileMetadata.TimeModified.IsZero() {
			if err := os.Chtimes(file.FileMetadata.Path, time.Now(), file.FileMetadata.TimeModified); err != nil {
				t.Errorf("Failed to change the access and modification times of the file: %v", err)
			}
		}
	}

	*nodes = append(*nodes, FileSystemNode{
		Path:        file.FileMetadata.Path,
		IsDirectory: file.FileMetadata.IsDirectory,
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

	files := CreateSortedFileSystemFiles(t, "", input)
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
			file.FileMetadata.Path = filepath.Join(directory, file.FileMetadata.Path)
			if i == 0 || file.FileMetadata.DirectoryPath != previousDirectoryPath {
				TestingCreateDirectoryAll(t, file.FileMetadata.DirectoryPath)
			}
			previousDirectoryPath = file.FileMetadata.DirectoryPath
			testingIfFileWriteItAndAppendFileSystemNode(t, file, &fileSystemNodes)
		}
	}

	return temporaryDirectories, fileSystemNodes
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
	files := CreateSortedFileSystemFiles(t, directory, input)

	for i := range files {
		if previousDirectoryPath != files[i].FileMetadata.DirectoryPath {
			TestingCreateDirectoryAll(t, files[i].FileMetadata.DirectoryPath)
			previousDirectoryPath = files[i].FileMetadata.DirectoryPath
		}
		testingIfFileWriteItAndAppendFileSystemNode(t, files[i], &nodes)
	}

	return directory, nodes
}
