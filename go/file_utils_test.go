package main

import (
	"os"
	"path/filepath"
	"sort"
	"testing"
	"time"
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

	if !dirDetail.IsDirectory {
		t.Errorf("Want IsDirectory to be true, got %v", dirDetail.IsDirectory)
	}

	if fileDetail.Path != filePath {
		t.Errorf("Want Path %v, got %v", filePath, fileDetail.Path)
	}

	if fileDetail.Size != int64(len(testText)) {
		t.Errorf("Want Size %v, got %v", len(testText), fileDetail.Size)
	}

	if fileDetail.IsDirectory {
		t.Errorf("Want IsDirectory to be false, got %v", fileDetail.IsDirectory)
	}

	if err == nil {
		t.Errorf("Want an error when trying to get details of a non-existent file, but got none")
	}
}

// TODO: WIP
func TestGetFilteredFileDetailsSliceFromDirectoryTree(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "testTempDir")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	subDir := filepath.Join(tmpDir, "subDir")
	file1 := filepath.Join(tmpDir, "file1.txt")
	file2 := filepath.Join(tmpDir, "file2.txt")
	file3 := filepath.Join(subDir, "file3.txt")

	os.Mkdir(subDir, 0755)
	os.WriteFile(file1, []byte("content"), 0666)
	os.WriteFile(file2, []byte(""), 0666)
	os.WriteFile(file3, []byte("content"), 0666)

	modTime := time.Now().Add(-1 * time.Hour)
	os.Chtimes(file1, modTime, modTime)
	os.Chtimes(file2, modTime, modTime)
	os.Chtimes(file3, modTime, modTime)

	tests := []struct {
		name           string
		rootFilePath   string
		fileFilterMode FileFilterMode
		want           []FileDetail
		wantErr        bool
	}{
		{
			name:           "FilesOnly",
			rootFilePath:   tmpDir,
			fileFilterMode: files,
			want: []FileDetail{
				{Path: filepath.Join(tmpDir, "file1.txt"), Size: 7, IsDirectory: false},
				{Path: filepath.Join(tmpDir, "file2.txt"), Size: 0, IsDirectory: false},
				{Path: filepath.Join(subDir, "file3.txt"), Size: 7, IsDirectory: false},
			},
			wantErr: false,
		},
		{
			name:           "DirectoriesOnly",
			rootFilePath:   tmpDir,
			fileFilterMode: directories,
			want: []FileDetail{
				{Path: tmpDir, IsDirectory: true},
				{Path: subDir, IsDirectory: true},
			},
			wantErr: false,
		},
		{
			name:           "FilesWithoutZeroByteFiles",
			rootFilePath:   tmpDir,
			fileFilterMode: filesWithoutZeroByteFiles,
			want: []FileDetail{
				{Path: filepath.Join(tmpDir, "file1.txt"), Size: 7, IsDirectory: false},
				{Path: filepath.Join(subDir, "file3.txt"), Size: 7, IsDirectory: false},
			},
			wantErr: false,
		},
		{
			name:           "InvalidRootPath",
			rootFilePath:   "/invalid/path",
			fileFilterMode: files,
			want:           nil,
			wantErr:        true,
		},
		// {
		// 	name:           "PermissionError",
		// 	rootFilePath:   tmpDir,
		// 	fileFilterMode: files,
		// 	want: []FileDetail{
		// 		{Path: filepath.Join(tmpDir, "file1.txt"), Size: 7, IsDirectory: false},
		// 		{Path: filepath.Join(tmpDir, "file2.txt"), Size: 0, IsDirectory: false},
		// 	},
		// 	wantErr: false,
		// },
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// if tt.name == "PermissionError" {
			// 	err := os.Chmod(subDir, 0000)
			// 	if err != nil {
			// 		t.Fatal(err)
			// 	}
			// 	defer func() {
			// 		if err := os.Chmod(subDir, 0755); err != nil {
			// 			t.Log("Warning: Failed to restore permissions for subDir:", err)
			// 		}
			// 	}()
			// }

			got, err := getFilteredFileDetailsSliceFromDirectoryTree(tt.rootFilePath, tt.fileFilterMode)
			if (err != nil) != tt.wantErr {
				t.Errorf("getFilteredFileDetailsSliceFromDirectoryTree() error: %v, wantErr %v", err, tt.wantErr)
				return
			}

			// If no error, check the content of the result
			if !tt.wantErr && err == nil {
				// Sort the slices by file path to ensure a consistent order for comparison.
				sort.Slice(got, func(i, j int) bool { return got[i].Path < got[j].Path })
				sort.Slice(tt.want, func(i, j int) bool { return tt.want[i].Path < tt.want[j].Path })

				// Compare the length first.
				if len(got) != len(tt.want) {
					t.Fatalf("got %d file details; want %d", len(got), len(tt.want))
				}

				// Now compare the fields, ignoring the ModificationTime.
				for i := range got {
					if got[i].Path != tt.want[i].Path ||
						got[i].Size != tt.want[i].Size ||
						got[i].IsDirectory != tt.want[i].IsDirectory {
						t.Errorf("getFilteredFileDetailsSliceFromDirectoryTree() details mismatch; got %+v, want %+v", got[i], tt.want[i])
					}
				}
			}
		})
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
