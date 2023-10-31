package internal

import (
	"path/filepath"
	"testing"
)

// TODO: check this test
func TestJoinOutputBasePathWithRelativeInputPath(t *testing.T) {
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
			inputBasePath:  "/home/user/source",
			inputFullPath:  "/home/user/source/subdir/file.txt",
			outputBasePath: "/home/user/dest",
			expected:       filepath.FromSlash("/home/user/dest/subdir/file.txt"),
			expectErr:      false,
		},
		{
			name:           "Error case",
			inputBasePath:  "/home/user/source",
			inputFullPath:  "/home/other/folder/file.txt",
			outputBasePath: "/home/user/dest",
			expected:       filepath.FromSlash("/home/other/folder/file.txt"),
			expectErr:      false,
		},
		{
			name:           "Equivalent paths",
			inputBasePath:  "/home/user/source",
			inputFullPath:  "/home/user/source",
			outputBasePath: "/home/user/dest",
			expected:       filepath.FromSlash("/home/user/dest"),
			expectErr:      false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := joinOutputBasePathWithRelativeInputPath(tt.inputBasePath, tt.inputFullPath, tt.outputBasePath)
			if (err != nil) != tt.expectErr {
				t.Fatalf("expected error: %v, got %v", tt.expectErr, err)
			}
			if err == nil && result != tt.expected {
				t.Fatalf("expected: %s, got %s", tt.expected, result)
			}
		})
	}
}
