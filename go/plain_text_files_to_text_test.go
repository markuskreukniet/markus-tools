package main

import (
	"fmt"
	"path/filepath"
	"strings"
	"testing"

	"github.com/markuskreukniet/markus-tools/go/utils/test"
)

func testingLastPathElementOnNewline(filePath string) string {
	return fmt.Sprintf("%s\n", filepath.Base(filePath))
}

func testingCreateContentString(filePath string, index int) string {
	return fmt.Sprintf("content %s %d 1\ncontent %s %d 2", filePath, index, filePath, index)
}

// TODO: naming
func testingTest(t *testing.T, directoryWithOptionalFileAsStrings []string, builder *strings.Builder) {
	t.Helper()
	test.TestingWriteString(t, directoryWithOptionalFileAsStrings[2]+"\n"+directoryWithOptionalFileAsStrings[3], builder)
}

func TestPlainTextFilesToText(t *testing.T) {
	// arrange
	testCases := []struct {
		metadata test.TestCaseMetadata
		input    string
	}{
		{
			metadata: test.TestingCreateTestCaseMetadataWithNameBasicAndWantErrFalse(),
			input: `
				,,txt 0.txt,content 0\ncontent 0;
				,,jpg 0.jpg,;
				empty,,,;
				directory 1/empty,,,;
				directory 1/directory 2,,txt 1-2.txt,content directory 1/directory 2 1-2\ncontent directory 1/directory 2 1-2;
				directory 1/directory 2,,txt 1-2 2.txt,content directory 1/directory 2 1-2 2\ncontent directory 1/directory 2 1-2 2;
			`,
		},
		{
			metadata: test.TestCaseMetadata{
				Name:    "Empty Input",
				WantErr: false,
			},
			input: "",
		},
	}

	// testCases
	for _, tc := range testCases {
		t.Run(tc.metadata.Name, func(t *testing.T) {
			// arrange and teardown
			directory, fileSystemNodes := test.TestingCreateFilesAndDirectories(t, tc.input)
			defer test.TestingRemoveDirectoryTree(t, directory)
			var builder strings.Builder
			if directory != "" {
				// TODO: duplicate
				directoriesWithOptionalFileAsDelimitedCommaString := strings.Split(strings.TrimSuffix(strings.TrimSpace(tc.input), ";"), ";")
				index := 0
				for i, delimitedCommaString := range directoriesWithOptionalFileAsDelimitedCommaString {
					// TODO: duplicate
					directoryWithOptionalFileAsStrings := strings.Split(strings.TrimSpace(delimitedCommaString), ",")
					if directoryWithOptionalFileAsStrings[3] != "" {
						testingTest(t, directoryWithOptionalFileAsStrings, &builder)
						index = i + 1
						break
					}
				}
				for i := index; i < len(directoriesWithOptionalFileAsDelimitedCommaString); i++ {
					directoryWithOptionalFileAsStrings := strings.Split(directoriesWithOptionalFileAsDelimitedCommaString[i], ",")
					if directoryWithOptionalFileAsStrings[3] != "" {
						test.TestingWriteString(t, "\n\n", &builder)
						testingTest(t, directoryWithOptionalFileAsStrings, &builder)
					}
				}
			}

			// act
			outcome, err := plainTextFilesToText(fileSystemNodes)

			// assert
			test.TestingAssertErrorToWantErrorAndOutcomeToBuilderString(t, err, tc.metadata.WantErr, builder, outcome)
		})
	}
}

// func TestPlainTextFilesToText(t *testing.T) {
// 	// arrange

// 	Directory1 = "directory 1"
// 	Directory2 = "directory 2"
// 	Directory2WithDirectory3 = filepath.Join(Directory2, "directory 2")
// 	Directory2WithDirectory4 = filepath.Join(Directory2, "directory 4")

// 	TxtFile1 = filepath.Join(Directory1, "file 1.txt")
// 	TxtFile3 = filepath.Join(Directory2WithDirectory3, "file 3.txt")
// 	TxtFile6 = filepath.Join(Directory2WithDirectory4, "file 6.txt")
// 	JpgFile4 = filepath.Join(Directory1, "file 4.jpg")

// 	//

// 	input := `
// 		directory 1,,txt 1.txt,content directory 1 1\ncontent directory 1 1;
// 		directory 1,,jpg 1.jpg,;
// 		directory 2/directory 2,,txt 2-3.txt,content directory 2/directory 2 2-3\ncontent directory 2/directory 2 2-3;
// 		directory 2/directory 2,,txt 2-3 2.txt,content directory 2/directory 2 2-3 2\ncontent directory 2/directory 2 2-3 2;
// 	`

// 	TestingCreateFilesAndDirectories

// 	// arrange
// 	fileSystemPathEndParts := test.FileSystemPathEndParts{
// 		DirectoryPathEndParts: []string{test.Directory1, test.Directory2WithDirectory3, test.Directory2WithDirectory4},
// 		FilePathEndParts:      []string{test.TxtFile1, test.TxtFile3, test.TxtFile6, test.JpgFile4},
// 	}
// 	plainTextFilePathEndParts := []string{test.TxtFile1, test.TxtFile3, test.TxtFile6}
// 	fileSystemNodes := []utils.FileSystemNode{
// 		{
// 			Path:        test.TxtFile1,
// 			IsDirectory: false,
// 		},
// 		{
// 			Path:        test.Directory2,
// 			IsDirectory: true,
// 		},
// 	}
// 	var emptyFileSystemNodes []utils.FileSystemNode

// 	testCases := []struct {
// 		metadata                  test.TestCaseMetadata
// 		plainTextFilePathEndParts []string
// 		fileSystemNodes           []utils.FileSystemNode
// 	}{
// 		{
// 			metadata:                  test.TestingCreateTestCaseMetadataWithNameBasicAndWantErrFalse(),
// 			plainTextFilePathEndParts: plainTextFilePathEndParts,
// 			fileSystemNodes:           fileSystemNodes,
// 		},
// 		{
// 			metadata:                  test.TestingCreateTestCaseMetadataWithNameEmptyFileSystemNodesAndWantErrFalse(),
// 			plainTextFilePathEndParts: test.EmptyPathEndParts,
// 			fileSystemNodes:           emptyFileSystemNodes,
// 		},
// 	}

// 	for _, tc := range testCases {
// 		t.Run(tc.metadata.Name, func(t *testing.T) {
// 			// arrange and teardown
// 			directory := test.TestingCreateTempFileSystemStructureOrGetEmptyString(t, fileSystemPathEndParts)
// 			defer test.TestingRemoveDirectoryTree(t, directory)
// 			for i := range tc.fileSystemNodes {
// 				tc.fileSystemNodes[i].Path = filepath.Join(directory, tc.fileSystemNodes[i].Path)
// 			}
// 			var builder strings.Builder
// 			// TODO: tc.PlainTextFilePathEndParts to full path happens also in testingCreateTempFileSystemStructureOrGetEmptyString?
// 			if len(tc.plainTextFilePathEndParts) > 0 {
// 				fullPath := filepath.Join(directory, tc.plainTextFilePathEndParts[0])
// 				content := testingCreateContentString(fullPath, 0)
// 				test.TestingWriteFileContent(t, fullPath, content)
// 				test.TestingWriteString(t, testingLastPathElementOnNewline(tc.plainTextFilePathEndParts[0]), &builder)
// 				test.TestingWriteString(t, content, &builder)
// 				for i := 1; i < len(tc.plainTextFilePathEndParts); i++ {
// 					fullPath := filepath.Join(directory, tc.plainTextFilePathEndParts[i])
// 					content := testingCreateContentString(fullPath, i)
// 					test.TestingWriteFileContent(t, fullPath, content)
// 					test.TestingWriteString(t, "\n\n", &builder)
// 					test.TestingWriteString(t, testingLastPathElementOnNewline(tc.plainTextFilePathEndParts[i]), &builder)
// 					test.TestingWriteString(t, content, &builder)
// 				}
// 			}

// 			// act
// 			outcome, err := plainTextFilesToText(tc.fileSystemNodes)

// 			// assert
// 			test.TestingAssertErrorToWantErrorAndOutcomeToBuilderString(t, err, tc.metadata.WantErr, builder, outcome)
// 		})
// 	}
// }
