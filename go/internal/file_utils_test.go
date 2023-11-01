package internal

import (
	"os"
	"path/filepath"
	"testing"
)

func fatalLogIfError(t *testing.T, err error) {
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
	fatalLogIfError(t, err)

	fileDetail, err := getFileDetail(filePath)
	fatalLogIfError(t, err)

	_, err = getFileDetail(nonExistentFilePath)

	// Assert
	if dirDetail.Path != tempDir {
		t.Errorf("Expected Path %v, got %v", tempDir, dirDetail.Path)
	}

	if !dirDetail.IsDirectory {
		t.Errorf("Expected IsDirectory to be true, got %v", dirDetail.IsDirectory)
	}

	if fileDetail.Path != filePath {
		t.Errorf("Expected Path %v, got %v", filePath, fileDetail.Path)
	}

	if fileDetail.Size != int64(len(testText)) {
		t.Errorf("Expected Size %v, got %v", len(testText), fileDetail.Size)
	}

	if fileDetail.IsDirectory {
		t.Errorf("Expected IsDirectory to be false, got %v", fileDetail.IsDirectory)
	}

	if err == nil {
		t.Errorf("Expected an error when trying to get details of a non-existent file, but got none")
	}
}

func TestJoinOutputBasePathWithRelativeInputPath(t *testing.T) {
	const inputBasePath string = "/home/user/source"
	const inputFullPath string = "/home/user/source/directory/file.txt"
	const outputBasePath string = "/home/user/destination"
	const joinedOutputBasePathWithRelativeInputPath string = "/home/user/destination/directory/file.txt"

	tests := []struct {
		name           string
		inputBasePath  string
		inputFullPath  string
		outputBasePath string
		expected       string
		expectErr      bool
	}{
		{
			name:           "Basic",
			inputBasePath:  inputBasePath,
			inputFullPath:  inputFullPath,
			outputBasePath: outputBasePath,
			expected:       filepath.FromSlash(joinedOutputBasePathWithRelativeInputPath),
			expectErr:      false,
		},
		{
			name:           "Empty inputBasePath",
			inputBasePath:  "",
			inputFullPath:  inputFullPath,
			outputBasePath: outputBasePath,
			expected:       "",
			expectErr:      true,
		},
		{
			name:           "Empty inputFullPath",
			inputBasePath:  inputBasePath,
			inputFullPath:  "",
			outputBasePath: outputBasePath,
			expected:       "",
			expectErr:      true,
		},
		{
			name:           "Equivalent input paths",
			inputBasePath:  inputBasePath,
			inputFullPath:  inputBasePath,
			outputBasePath: outputBasePath,
			expected:       filepath.FromSlash(outputBasePath),
			expectErr:      false,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result, err := joinOutputBasePathWithRelativeInputPath(test.inputBasePath, test.inputFullPath, test.outputBasePath)
			if (err != nil) != test.expectErr {
				t.Fatalf("expected error: %v, got %v", test.expectErr, err)
			}
			if err == nil && result != test.expected {
				t.Fatalf("expected: %s, got %s", test.expected, result)
			}
		})
	}
}
