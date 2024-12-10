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
			AbsoluteDirectoryPath: directoryPath, // TODO: here it is not an absolute path, also fix in Kotlin?
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

	directoryPath := tMustCreateTemporaryDirectory(t) // TODO: also do it like this in Kotlin code
	files := createFilesData(t, directoryPath, input)

	createdDirectoryPaths := make(map[string]struct{})

	for _, file := range files {
		if _, exists := createdDirectoryPaths[file.CompleteFileInfo.AbsoluteDirectoryPath]; !exists {
			tMustCreateDirectoryAll(t, file.CompleteFileInfo.AbsoluteDirectoryPath)
			createdDirectoryPaths[file.CompleteFileInfo.AbsoluteDirectoryPath] = struct{}{}
		}
		ifFileThenWriteAndChangeTimes(t, file)
	}

	return directoryPath
}

// func WriteFilesByMultipleInputs(t *testing.T, input string) ([]string, []FileSystemNode) {
// 	if IsBlank(input) {
// 		return nil, nil
// 	}

// 	files := createFilesData(t, "", input)

// 	if len(files) == 0 {
// 		return nil, nil
// 	}

// 	sort.Slice(files, func(i, j int) bool {
// 		return files[i].CompleteFileInfo.AbsolutePath < files[j].CompleteFileInfo.AbsolutePath
// 	})

// 	previousDirectoryPath := toRootDirectoryPath(files[0].CompleteFileInfo.AbsoluteDirectoryPath)
// 	fileGroups := [][]FileData{{files[0]}}
// 	index := 0

// 	for _, file := range files[1:] {

// 	}

// 	previousDirectoryPath = ""
// 	var temporaryDirectoryPaths []string
// 	var fileSystemNodes []FileSystemNode

// 	return temporaryDirectoryPaths, fileSystemNodes
// }

func tMustCreateTemporaryDirectory(t *testing.T) string {
	result, err := os.MkdirTemp("", "markus-tools go test")
	return TMust(t, result, err)
}

// TODO: should receive FileData?
func tMustWriteFile(t *testing.T, filePath string, content string) {
	TMustErr(t, os.WriteFile(filePath, []byte(content), 0666))
}

func tMustCreateDirectoryAll(t *testing.T, filePath string) {
	TMustErr(t, os.MkdirAll(filePath, 0755))
}

func changeFileTimes(file FileData) error {
	return os.Chtimes(file.CompleteFileInfo.AbsolutePath, time.Now(), file.CompleteFileInfo.TimeModified)
}

func tMustChangeFileTimes(t *testing.T, file FileData) {
	TMustErr(t, changeFileTimes(file))
}

func ifFileThenWriteAndChangeTimes(t *testing.T, file FileData) {
	if !file.CompleteFileInfo.IsDirectory {
		tMustWriteFile(t, file.CompleteFileInfo.AbsolutePath, file.Content)
		tMustChangeFileTimes(t, file)
	}
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

func WriteFileWithContentAndIndex(t *testing.T, filePath string, index int) string {
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
