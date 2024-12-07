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

func createFileData(directoryPath, inputLine string) (FileData, error) {
	fields := strings.Split(inputLine, ",")
	directoryPath = filepath.Join(directoryPath, filepath.FromSlash(fields[0]))
	content := fields[3]
	name := fields[2]
	filePath := filepath.Join(directoryPath, name)
	isDirectory := name == ""

	var timeModified time.Time
	if fields[1] != "" {
		var err error
		timeModified, err = time.Parse(time.RFC3339, fields[1])
		if err != nil {
			return FileData{}, err
		}
	}

	return FileData{
		Content: content,
		CompleteFileInfo: CompleteFileInfo{
			Name:                  name,
			AbsoluteDirectoryPath: directoryPath,
			AbsolutePath:          filePath,
			TimeModified:          timeModified,
			Size:                  0, // TODO: convert content to size?
			IsDirectory:           isDirectory,
		},
	}, nil
}

func tMustCreateFileData(t *testing.T, directoryPath, inputLine string) FileData {
	result, err := createFileData(directoryPath, inputLine)
	return TMust(t, result, err)
}

func createFilesDataWithEmptyDirectoryPath(t *testing.T, rawDelimitedSemicolonString string) []FileData {
	return createFilesData(t, "", rawDelimitedSemicolonString)
}

func createFilesData(t *testing.T, directoryPath, rawDelimitedSemicolonString string) []FileData {
	var files []FileData
	var inputLine []rune
	isCreatingInputLine := false
	rawDelimitedSemicolonString = strings.TrimSpace(rawDelimitedSemicolonString)

	for _, r := range rawDelimitedSemicolonString {
		if isCreatingInputLine {
			if r != ';' {
				inputLine = append(inputLine, r)
			} else {
				files = append(files, tMustCreateFileData(t, directoryPath, string(inputLine)))
				inputLine = nil
				isCreatingInputLine = false
			}
		} else if !unicode.IsSpace(r) {
			inputLine = append(inputLine, r)
			isCreatingInputLine = true
		}
	}

	return files
}

func WriteFilesBySingleInput(t *testing.T, input string) string {
	if IsBlank(input) {
		return ""
	}

	files := createFilesDataWithEmptyDirectoryPath(t, input)

	if len(files) == 0 {
		return ""
	}

	directoryPath := tMustCreateTemporaryDirectory(t)

	for _, file := range files {
		joinAbsolutePaths(directoryPath, &file)
		if !file.CompleteFileInfo.IsDirectory {
			tMustWriteFile(t, file.CompleteFileInfo.AbsolutePath, file.Content)
		}
	}

	return directoryPath
}

func tMustCreateTemporaryDirectory(t *testing.T) string {
	result, err := os.MkdirTemp("", "markus-tools go test")
	return TMust(t, result, err)
}

func joinAbsolutePaths(directoryPath string, file *FileData) {
	file.CompleteFileInfo.AbsoluteDirectoryPath =
		filepath.Join(directoryPath, file.CompleteFileInfo.AbsoluteDirectoryPath)
	file.CompleteFileInfo.AbsolutePath = filepath.Join(directoryPath, file.CompleteFileInfo.AbsolutePath)
}

// TODO: should receive FileData?
func tMustWriteFile(t *testing.T, filePath string, content string) {
	TMustErr(t, os.WriteFile(filePath, []byte(content), 0666))
}

func tMustCreateDirectoryAll(t *testing.T, filePath string) {
	TMustErr(t, os.MkdirAll(filePath, 0755))
}

// old
func createFileSystemFileByInputLine(t *testing.T, directoryPath, inputLine string) FileSystemFile {
	t.Helper()

	fields := strings.Split(inputLine, ",")
	directoryPath = filepath.Join(directoryPath, filepath.FromSlash(fields[0]))
	data := fields[3]
	name := fields[2]
	filePath := filepath.Join(directoryPath, name)
	isDirectory := name == ""

	var timeModified time.Time
	if fields[1] != "" {
		var err error
		timeModified, err = time.Parse(time.RFC3339, fields[1])
		if err != nil {
			t.Errorf("Failed to parse time: %v", err)
		}
	}

	return CreateFileSystemFile(data, CreateFileMetadata(name, directoryPath, filePath, "", timeModified, 0, isDirectory))
}

// old
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

// TODO: wrong naming
func TestingWriteFileWithContentAndIndex(t *testing.T, filePath string, index int) string {
	t.Helper()
	writtenContent := fmt.Sprintf("content %d", index)
	tMustWriteFile(t, filePath, writtenContent)
	return writtenContent
}

func testingIfFileWriteItAndAppendFileSystemNode(t *testing.T, file FileSystemFile, nodes *[]FileSystemNode) {
	t.Helper()

	if !file.FileMetadata.IsDirectory {
		tMustWriteFile(t, file.FileMetadata.Path, file.Data)
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

	if IsBlank(input) {
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

	for _, file := range files[1:] {
		rootDirectoryPath := toRootDirectoryPath(file.FileMetadata.DirectoryPath)
		if rootDirectoryPath == "." || rootDirectoryPath != previousRootDirectoryPath {
			fileGroups = append(fileGroups, []FileSystemFile{file})
			previousRootDirectoryPath = rootDirectoryPath
			index++
		} else {
			fileGroups[index] = append(fileGroups[index], file)
		}
	}

	var temporaryDirectories []string
	var fileSystemNodes []FileSystemNode
	var previousDirectoryPath string

	for _, group := range fileGroups {
		directory := tMustCreateTemporaryDirectory(t)
		temporaryDirectories = append(temporaryDirectories, directory)
		for i, file := range group {
			file.FileMetadata.DirectoryPath = filepath.Join(directory, file.FileMetadata.DirectoryPath)
			file.FileMetadata.Path = filepath.Join(directory, file.FileMetadata.Path)
			if i == 0 || file.FileMetadata.DirectoryPath != previousDirectoryPath {
				tMustCreateDirectoryAll(t, file.FileMetadata.DirectoryPath)
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

	if IsBlank(input) {
		return "", nil
	}

	var nodes []FileSystemNode
	var previousDirectoryPath string
	directory := tMustCreateTemporaryDirectory(t)
	files := CreateSortedFileSystemFiles(t, directory, input)

	for i := range files {
		if previousDirectoryPath != files[i].FileMetadata.DirectoryPath {
			tMustCreateDirectoryAll(t, files[i].FileMetadata.DirectoryPath)
			previousDirectoryPath = files[i].FileMetadata.DirectoryPath
		}
		testingIfFileWriteItAndAppendFileSystemNode(t, files[i], &nodes)
	}

	return directory, nodes
}
