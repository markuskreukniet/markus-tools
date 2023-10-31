package internal

import (
	"path/filepath"
	"testing"
)

func TestJoinOutputBasePathWithRelativeInputPath(t *testing.T) {
	const inputBasePath string = "/home/user/source"
	const inputFullPath string = "/home/user/source/directory/file.txt"
	const outputBasePath string = "/home/user/destination"
	const otherFullPath string = "/home/other/folder/file.txt"
	const joinedOutputBasePathWithRelativeInputPath string = "/home/user/destination/directory/file.txt"

	tests := []struct {
		name           string
		inputBasePath  string
		inputFullPath  string
		outputBasePath string
		expected       string
		expectedErr    bool
	}{
		{
			name:           "Basic",
			inputBasePath:  inputBasePath,
			inputFullPath:  inputFullPath,
			outputBasePath: outputBasePath,
			expected:       filepath.FromSlash(joinedOutputBasePathWithRelativeInputPath),
			expectedErr:    false,
		},
		{
			name:           "Empty inputBasePath",
			inputBasePath:  "",
			inputFullPath:  inputFullPath,
			outputBasePath: outputBasePath,
			expected:       "",
			expectedErr:    true,
		},
		{
			name:           "Empty inputFullPath",
			inputBasePath:  inputBasePath,
			inputFullPath:  "",
			outputBasePath: outputBasePath,
			expected:       "",
			expectedErr:    true,
		},
		{
			name:           "Equivalent input paths",
			inputBasePath:  inputBasePath,
			inputFullPath:  inputBasePath,
			outputBasePath: outputBasePath,
			expected:       filepath.FromSlash(outputBasePath),
			expectedErr:    false,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result, err := joinOutputBasePathWithRelativeInputPath(test.inputBasePath, test.inputFullPath, test.outputBasePath)
			if (err != nil) != test.expectedErr {
				t.Fatalf("expected error: %v, got %v", test.expectedErr, err)
			}
			if err == nil && result != test.expected {
				t.Fatalf("expected: %s, got %s", test.expected, result)
			}
		})
	}
}
