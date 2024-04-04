package main

import (
	"bufio"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/markuskreukniet/markus-tools/go/utils"
	"github.com/markuskreukniet/markus-tools/go/utils/test"
)

type filePathEndPartContent struct {
	filePathEndPart string
	content         string
}

func testingHaveDirectoryTreesSameFilePaths(t *testing.T, sourceDirectory, destinationDirectory string) bool {
	t.Helper()
	if sourceDirectory == "" || destinationDirectory == "" {
		return false
	}

	// Do the directory trees have the same file paths?
	haveSameFilePaths := true
	err := filepath.Walk(sourceDirectory, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		relativePath, err := filepath.Rel(sourceDirectory, path)
		if err != nil {
			return err
		}
		destinationPath := filepath.Join(destinationDirectory, relativePath)
		exists, err := utils.FileOrDirectoryExists(destinationPath)
		if err != nil {
			return err
		}
		if !exists {
			haveSameFilePaths = false
		}
		return nil
	})
	if err != nil {
		t.Errorf("Failed to walk through directory: %v", err)
	}
	return haveSameFilePaths
}

func testingContainsTxtFile4(stringSlice []string) bool {
	for _, item := range stringSlice {
		if item == test.TxtFile4 {
			return true
		}
	}
	return false
}

func TestSynchronizeDirectoryTrees2(t *testing.T) {
	// arrange
	newContent := "content directory 2/directory\ncontent 3 2-3 2 new"
	input := `
		empty,,,;
		,,txt 0.txt,;
		directory 1,,txt 1.txt,;
		directory 2/directory 3,,txt 2-3.txt,;
	`
	sourceInput := input + `
		,,jpg 0.jpg,;
		directory 2/empty,,,;
		directory 2/directory 3,2020-02-20T20:40:41Z,txt 2-3 2.txt,` + newContent + ";"
	destinationInput := input + `
		directory 2/directory 3/empty,,,;
		directory 2/directory 3,2020-02-20T20:40:40Z,txt 2-3 2.txt,content directory 2/directory\ncontent 3 2-3 2 old;
		directory 2/directory 3,,txt 2-3 3.txt,;
	`

	// create testCases
	testCases := []struct {
		metadata          test.TestCaseMetadata
		sourceInput       string
		destinationInput  string
		updatedFile       filePathEndPartContent
		wantSameFilePaths bool
	}{
		{
			metadata:         test.TestingCreateTestCaseMetadataWithNameBasicAndWantErrFalse(),
			sourceInput:      sourceInput,
			destinationInput: destinationInput,
			updatedFile: filePathEndPartContent{
				filePathEndPart: "directory 2/directory 3/txt 2-3 2.txt", // TODO: should be a combination of two vars and "/"?
				content:         newContent,
			},
			wantSameFilePaths: true,
		},
		{
			metadata:          test.TestingCreateTestCaseMetadataWithWantErrTrue("Empty destinationInput"),
			sourceInput:       sourceInput,
			destinationInput:  "",
			updatedFile:       filePathEndPartContent{},
			wantSameFilePaths: false,
		},
	}

	// run testCases
	for _, tc := range testCases {
		t.Run(tc.metadata.Name, func(t *testing.T) {
			// arrange and teardown
			sourceDirectory, _ := test.TestingCreateFilesAndDirectories(t, tc.sourceInput)
			defer test.TestingRemoveDirectoryTree(t, sourceDirectory)
			destinationDirectory, _ := test.TestingCreateFilesAndDirectories(t, tc.destinationInput)
			defer test.TestingRemoveDirectoryTree(t, destinationDirectory)

			// act
			err := synchronizeDirectoryTrees(sourceDirectory, destinationDirectory)

			// assert
			test.TestingAssertErrorToWantError(t, err, tc.metadata.WantErr)
			haveSameFilePaths := testingHaveDirectoryTreesSameFilePaths(t, sourceDirectory, destinationDirectory)
			if tc.wantSameFilePaths && !haveSameFilePaths {
				t.Errorf("The source and destination directory trees do not have the same file paths.")
			}
			haveSameFilePaths = testingHaveDirectoryTreesSameFilePaths(t, destinationDirectory, sourceDirectory)
			if tc.wantSameFilePaths && !haveSameFilePaths {
				t.Errorf("The destination and source directory trees do not have the same file paths.")
			}
			if tc.updatedFile.filePathEndPart != "" && tc.updatedFile.content != "" {
				content, err := os.ReadFile(test.ToFilePathFromSlashAndJoin(destinationDirectory, tc.updatedFile.filePathEndPart))
				if err != nil {
					t.Fatalf("Failed to read file: %s", err)
				}
				test.TestingAssertEqualStrings(t, string(content), tc.updatedFile.content)
			}
		})
	}
}

func TestSynchronizeDirectoryTrees(t *testing.T) {
	// arrange
	sourceFileSystemPathEndParts := test.FileSystemPathEndParts{
		DirectoryPathEndParts: []string{test.DirectoryEmpty, test.Directory1, test.Directory2WithDirectoryEmpty, test.Directory2WithDirectory3},
		FilePathEndParts:      []string{test.TxtFile1, test.TxtFile4, test.TxtFile5},
	}
	destinationFileSystemPathEndParts := test.FileSystemPathEndParts{
		DirectoryPathEndParts: []string{test.DirectoryEmpty, test.Directory2WithDirectory3},
		FilePathEndParts:      []string{test.TxtFile3, test.TxtFile4},
	}

	testCases := []struct {
		metadata                          test.TestCaseMetadata
		sourceFileSystemPathEndParts      test.FileSystemPathEndParts
		destinationFileSystemPathEndParts test.FileSystemPathEndParts
		wantSameFilePaths                 bool
	}{
		{
			metadata:                          test.TestingCreateTestCaseMetadataWithNameBasicAndWantErrFalse(),
			sourceFileSystemPathEndParts:      sourceFileSystemPathEndParts,
			destinationFileSystemPathEndParts: destinationFileSystemPathEndParts,
			wantSameFilePaths:                 true,
		},
		{
			metadata:                          test.TestingCreateTestCaseMetadataWithWantErrTrue("Empty DestinationPathEndParts"),
			sourceFileSystemPathEndParts:      sourceFileSystemPathEndParts,
			destinationFileSystemPathEndParts: test.EmptyFileSystemPathEndParts,
			wantSameFilePaths:                 false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.metadata.Name, func(t *testing.T) {
			// arrange and teardown
			sourceDirectory := test.TestingCreateTempFileSystemStructureOrGetEmptyString(t, tc.sourceFileSystemPathEndParts)
			defer test.TestingRemoveDirectoryTree(t, sourceDirectory)
			destinationDirectory := test.TestingCreateTempFileSystemStructureOrGetEmptyString(t, tc.destinationFileSystemPathEndParts)
			defer test.TestingRemoveDirectoryTree(t, destinationDirectory)

			// Some file systems have a resolution of one second, so we must wait a second.
			filePathTxtFile4 := ""
			writtenContent := ""
			if sourceDirectory != "" && destinationDirectory != "" && testingContainsTxtFile4(tc.sourceFileSystemPathEndParts.FilePathEndParts) && testingContainsTxtFile4(tc.destinationFileSystemPathEndParts.FilePathEndParts) {
				filePathTxtFile4 = filepath.Join(destinationDirectory, test.TxtFile4)
				test.TestingWriteFileContentWithContentAndIndex(t, filePathTxtFile4, 1)
				time.Sleep(time.Second)
				writtenContent = test.TestingWriteFileContentWithContentAndIndex(t, filepath.Join(sourceDirectory, test.TxtFile4), 2)
			}

			// act
			err := synchronizeDirectoryTrees(sourceDirectory, destinationDirectory)

			// assert
			test.TestingAssertErrorToWantError(t, err, tc.metadata.WantErr)
			haveSameFilePaths := testingHaveDirectoryTreesSameFilePaths(t, sourceDirectory, destinationDirectory)
			if tc.wantSameFilePaths && !haveSameFilePaths {
				t.Errorf("The source and destination directory trees do not have the same file paths.")
			}
			haveSameFilePaths = testingHaveDirectoryTreesSameFilePaths(t, destinationDirectory, sourceDirectory)
			if tc.wantSameFilePaths && !haveSameFilePaths {
				t.Errorf("The destination and source directory trees do not have the same file paths.")
			}
			if filePathTxtFile4 != "" {
				file, err := os.Open(filePathTxtFile4)
				if err != nil {
					t.Errorf("Failed to open file: %v", err)
				}
				defer file.Close()
				var builder strings.Builder
				scanner := bufio.NewScanner(file)
				for scanner.Scan() {
					test.TestingWriteString(t, scanner.Text(), &builder)
				}
				// TODO: if?
				err = scanner.Err()
				if err != nil {
					t.Errorf("Failed to read file content: %v", err)
				}
				test.TestingAssertEqualStrings(t, builder.String(), writtenContent)
			}
		})
	}
}

// TODO: use vars from arrange utils?
func TestJoinOutputBasePathWithRelativeInputPath(t *testing.T) {
	const inputBasePath string = "/home/user/source"
	const inputFullPath string = "/home/user/source/directory/file.txt"
	const outputBasePath string = "/home/user/destination"
	const joinedOutputBasePathWithRelativeInputPath string = "/home/user/destination/directory/file.txt"

	testCases := []struct {
		metadata       test.TestCaseMetadata
		inputBasePath  string
		inputFullPath  string
		outputBasePath string
		want           string
	}{
		{
			metadata:       test.TestingCreateTestCaseMetadataWithNameBasicAndWantErrFalse(),
			inputBasePath:  inputBasePath,
			inputFullPath:  inputFullPath,
			outputBasePath: outputBasePath,
			want:           filepath.FromSlash(joinedOutputBasePathWithRelativeInputPath),
		},
		{
			metadata:       test.TestingCreateTestCaseMetadataWithWantErrTrue("Empty InputBasePath"),
			inputBasePath:  "",
			inputFullPath:  inputFullPath,
			outputBasePath: outputBasePath,
			want:           "",
		},
		{
			metadata:       test.TestingCreateTestCaseMetadataWithWantErrTrue("Empty InputFullPath"),
			inputBasePath:  inputBasePath,
			inputFullPath:  "",
			outputBasePath: outputBasePath,
			want:           "",
		},
		{
			metadata:       test.TestingCreateTestCaseMetadata("Equivalent Input Paths", false),
			inputBasePath:  inputBasePath,
			inputFullPath:  inputBasePath,
			outputBasePath: outputBasePath,
			want:           filepath.FromSlash(outputBasePath),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.metadata.Name, func(t *testing.T) {
			result, err := joinOutputBasePathWithRelativeInputPath(tc.inputBasePath, tc.inputFullPath, tc.outputBasePath)
			test.TestingAssertErrorToWantError(t, err, tc.metadata.WantErr)
			if err == nil {
				test.TestingAssertEqualStrings(t, tc.want, result)
			}
		})
	}
}
