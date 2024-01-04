package main

import (
	"testing"
)

// TODO: there are duplicate things, such as statements and structs, probably also in other tests
func TestPlainTextFilesToText(t *testing.T) {
	// arrange
	// TODO: add non text files to filePathEndParts
	// directoryPathEndParts := []string{directory1, directory2WithDirectory3, directory2WithDirectory4}
	// filePathEndParts := []string{txtFile1, txtFile3, txtFile6}
	// fileSystemNodes := []FileSystemNode{
	// 	{
	// 		Path:        txtFile1,
	// 		IsDirectory: false,
	// 	},
	// 	{
	// 		Path:        directory2,
	// 		IsDirectory: true,
	// 	},
	// }

	// testCases := []struct {
	// 	Name                  string
	// 	DirectoryPathEndParts []string
	// 	FilePathEndParts      []string
	// 	WantErr               bool
	// }{
	// 	{
	// 		Name:                  "Basic",
	// 		DirectoryPathEndParts: directoryPathEndParts,
	// 		FilePathEndParts:      filePathEndParts,
	// 		WantErr:               false,
	// 	},
	// }

	// for _, tc := range testCases {
	// 	t.Run(tc.Name, func(t *testing.T) {
	// 		// arrange and tear down
	// 		directory, err := testingCreateTempFileSystemStructureOrGetEmptyString(tc.DirectoryPathEndParts, tc.FilePathEndParts)
	// 		if err != nil {
	// 			t.Fatalf("Failed to create the temporary directory: %v", err)
	// 		}
	// 		defer func() {
	// 			if err := os.RemoveAll(directory); err != nil {
	// 				t.Errorf("Failed to remove the temporary directory: %v", err)
	// 			}
	// 		}()
	// 		var builder strings.Builder
	// 		if len(fileSystemNodes) > 0 {
	// 			content := fmt.Sprintf("content %s %d\ncontent %s %d", fileSystemNodes[0].Path, 0)
	// 			testingWriteFileContent(t, fileSystemNodes[0].Path, content)
	// 			testingWriteString(t, content, &builder)
	// 		}
	// 		for _, node := range fileSystemNodes[1:] {
	// 			// TODO: duplicate
	// 			content := fmt.Sprintf("content %s %d\ncontent %s %d", fileSystemNodes[0].Path, 0)
	// 			testingWriteFileContent(t, fileSystemNodes[0].Path, content)
	// 			testingWriteString(t, content, &builder)
	// 		}

	// 		// act
	// 		text, err := plainTextFilesToText(fileSystemNodes)

	// 		// assert

	// 	})
	// }
}
