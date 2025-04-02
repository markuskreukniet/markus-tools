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

// TODO: should be tCreateFilesData? Look for same problem on other places
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

	directoryPath := tMustCreateTemporaryDirectory(t)
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

func WriteFilesByMultipleInputs(t *testing.T, input string) ([]string, []FileSystemNode) {
	if IsBlank(input) {
		return nil, nil
	}

	files := createFilesData(t, "", input)

	if len(files) == 0 {
		return nil, nil
	}

	sort.Slice(files, func(i, j int) bool {
		return files[i].CompleteFileInfo.AbsolutePath < files[j].CompleteFileInfo.AbsolutePath
	})

	previousSegment := getFirstPathSegment(files[0].CompleteFileInfo.AbsoluteDirectoryPath)
	fileGroups := [][]FileData{{files[0]}}
	index := 0

	for _, file := range files[1:] {
		currentSegment := getFirstPathSegment(file.CompleteFileInfo.AbsoluteDirectoryPath)
		if currentSegment == "." || currentSegment != previousSegment {
			fileGroups = append(fileGroups, []FileData{file})
			previousSegment = currentSegment
			index++
		} else {
			fileGroups[index] = append(fileGroups[index], file)
		}
	}

	// previousSegment = "" is unnecessary since the possible coming assignments are temporary directory paths,
	// which they were not before.
	var temporaryDirectoryPaths []string
	var fileSystemNodes []FileSystemNode

	for _, group := range fileGroups {
		directoryPath := tMustCreateTemporaryDirectory(t)
		temporaryDirectoryPaths = append(temporaryDirectoryPaths, directoryPath)
		for _, file := range group {
			file.CompleteFileInfo.AbsoluteDirectoryPath = filepath.Join(directoryPath, file.CompleteFileInfo.AbsoluteDirectoryPath)
			file.CompleteFileInfo.AbsolutePath = filepath.Join(directoryPath, file.CompleteFileInfo.AbsolutePath)
			if file.CompleteFileInfo.AbsoluteDirectoryPath != previousSegment {
				tMustCreateDirectoryAll(t, file.CompleteFileInfo.AbsoluteDirectoryPath)
			}
			previousSegment = file.CompleteFileInfo.AbsoluteDirectoryPath
			ifFileThenWriteAndChangeTimes(t, file)
			fileSystemNodes = append(fileSystemNodes, FileSystemNode{
				Path:        file.CompleteFileInfo.AbsolutePath,
				IsDirectory: file.CompleteFileInfo.IsDirectory,
			})
		}

	}

	return temporaryDirectoryPaths, fileSystemNodes
}

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

// TODO: obsolete?
func WriteFileWithContentAndIndex(t *testing.T, filePath string, index int) string {
	writtenContent := fmt.Sprintf("content %d", index)
	tMustWriteFile(t, filePath, writtenContent)
	return writtenContent
}

func getFirstPathSegment(filePath string) string {
	for {
		parentPath := filepath.Dir(filePath)
		if parentPath == "." || parentPath == filePath {
			break
		}
		filePath = parentPath
	}

	return filePath
}
