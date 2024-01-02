package main

import (
	"os"
	"path/filepath"
	"testing"
)

func testingFatalLogIfError(t *testing.T, err error) {
	t.Helper()
	if err != nil {
		t.Fatalf("getFileDetail() error: %v", err)
	}
}

func TestGetFileDetail(t *testing.T) {
	const testText string = "test text"

	// Arrange
	tempDir, err := os.MkdirTemp("", "testTempDir")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	filePath := filepath.Join(tempDir, "testFile.txt")
	err = os.WriteFile(filePath, []byte(testText), 0666)
	if err != nil {
		t.Fatalf("Failed to write to temp file: %v", err)
	}

	nonExistentFilePath := filepath.Join(tempDir, "nonExistentFile.txt")

	// Act
	dirDetail, err := getFileDetail(tempDir)
	testingFatalLogIfError(t, err)

	fileDetail, err := getFileDetail(filePath)
	testingFatalLogIfError(t, err)

	_, err = getFileDetail(nonExistentFilePath)

	// Assert
	if dirDetail.Path != tempDir {
		t.Errorf("Want Path %v, got %v", tempDir, dirDetail.Path)
	}

	if fileDetail.Path != filePath {
		t.Errorf("Want Path %v, got %v", filePath, fileDetail.Path)
	}

	if fileDetail.Size != int64(len(testText)) {
		t.Errorf("Want Size %v, got %v", len(testText), fileDetail.Size)
	}

	if err == nil {
		t.Errorf("Want an error when trying to get details of a non-existent file, but got none")
	}
}

func TestJoinOutputBasePathWithRelativeInputPath(t *testing.T) {
	const inputBasePath string = "/home/user/source"
	const inputFullPath string = "/home/user/source/directory/file.txt"
	const outputBasePath string = "/home/user/destination"
	const joinedOutputBasePathWithRelativeInputPath string = "/home/user/destination/directory/file.txt"

	testCases := []struct {
		Name           string
		InputBasePath  string
		InputFullPath  string
		OutputBasePath string
		Want           string
		WantErr        bool
	}{
		{
			Name:           "Basic",
			InputBasePath:  inputBasePath,
			InputFullPath:  inputFullPath,
			OutputBasePath: outputBasePath,
			Want:           filepath.FromSlash(joinedOutputBasePathWithRelativeInputPath),
			WantErr:        false,
		},
		{
			Name:           "Empty inputBasePath",
			InputBasePath:  "",
			InputFullPath:  inputFullPath,
			OutputBasePath: outputBasePath,
			Want:           "",
			WantErr:        true,
		},
		{
			Name:           "Empty inputFullPath",
			InputBasePath:  inputBasePath,
			InputFullPath:  "",
			OutputBasePath: outputBasePath,
			Want:           "",
			WantErr:        true,
		},
		{
			Name:           "Equivalent Input Paths",
			InputBasePath:  inputBasePath,
			InputFullPath:  inputBasePath,
			OutputBasePath: outputBasePath,
			Want:           filepath.FromSlash(outputBasePath),
			WantErr:        false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			result, err := joinOutputBasePathWithRelativeInputPath(tc.InputBasePath, tc.InputFullPath, tc.OutputBasePath)
			if (err != nil) != tc.WantErr {
				t.Fatalf("want error: %v, got %v", tc.WantErr, err)
			}
			if err == nil && result != tc.Want {
				t.Fatalf("want: %s, got %s", tc.Want, result)
			}
		})
	}
}
