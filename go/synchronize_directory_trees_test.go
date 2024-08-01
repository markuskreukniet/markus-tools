package main

import (
	"path/filepath"
	"testing"

	"github.com/markuskreukniet/markus-tools/go/utils"
)

type filePathEndPartContent struct {
	filePathEndPart string
	content         string
}

func TestSynchronizeDirectoryTrees(t *testing.T) {
	// arrange
	newContent := "content directory 2/directory\ncontent 3 2-3 2 new"
	input := `
		empty,,,;
		,,txt 0.txt,;
		directory 1,,txt 1.txt,;
		directory 2/directory 3,,txt 2-3.txt,;
	`

	// Some file systems have a resolution of one second, so the new file should be at least a second newer.
	sourceInput := input + `
		,,jpg 0.jpg,;
		directory 2/empty,,,;
		directory 2/directory 3,2020-02-20T20:40:41Z,txt 2-3 2.txt,` + newContent + ";"
	destinationInput := input + `
		directory 2/directory 3/empty,,,;
		directory 2/directory 3,2020-02-20T20:40:40Z,txt 2-3 2.txt,content directory 2/directory\ncontent 3 2-3 2 old;
		directory 2/directory 3,,txt 2-3 3.txt,;
	`

	// TODO: "Empty Input" is missing
	testCases := []utils.TestCaseBasicDoubleInput{
		utils.CreateTestCaseBasicDoubleInput(utils.CreateTestCaseBasic("Basic", sourceInput, sourceInput, false), destinationInput),
	}

	// run testCases
	for _, tc := range testCases {
		t.Run(tc.TestCaseBasic.Name, func(t *testing.T) {
			// arrange and teardown
			sourceDirectory, _ := utils.TestingWriteFilesByOneInput(t, tc.TestCaseBasic.Input)
			defer utils.TestingRemoveDirectoryTree(t, sourceDirectory)
			destinationDirectory, _ := utils.TestingWriteFilesByOneInput(t, tc.SecondInput)
			defer utils.TestingRemoveDirectoryTree(t, destinationDirectory)

			// act
			err := synchronizeDirectoryTrees(sourceDirectory, destinationDirectory)

			// assert
			utils.TestingAssertErrorToWantError(t, err, tc.TestCaseBasic.WantErr)

			areIdentical, err := utils.AreFileTreeDescendantsIdentical(sourceDirectory, destinationDirectory)
			if err != nil {
				t.Errorf("AreFileTreeDescendantsIdentical error: %v", err)
			}

			if !areIdentical {
				// TODO: message // TODO copied
				t.Errorf("outcome and wanted outcome are different")
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
		metadata       utils.TestCaseMetadata
		inputBasePath  string
		inputFullPath  string
		outputBasePath string
		want           string
	}{
		{
			metadata:       utils.CreateTestCaseMetadataWithNameBasicAndWantErrFalse(),
			inputBasePath:  inputBasePath,
			inputFullPath:  inputFullPath,
			outputBasePath: outputBasePath,
			want:           filepath.FromSlash(joinedOutputBasePathWithRelativeInputPath),
		},
		{
			metadata:       utils.CreateTestCaseMetadataWithWantErrTrue("Empty InputBasePath"),
			inputBasePath:  "",
			inputFullPath:  inputFullPath,
			outputBasePath: outputBasePath,
			want:           "",
		},
		{
			metadata:       utils.CreateTestCaseMetadataWithWantErrTrue("Empty InputFullPath"),
			inputBasePath:  inputBasePath,
			inputFullPath:  "",
			outputBasePath: outputBasePath,
			want:           "",
		},
		{
			metadata:       utils.CreateTestCaseMetadata("Equivalent Input Paths", false),
			inputBasePath:  inputBasePath,
			inputFullPath:  inputBasePath,
			outputBasePath: outputBasePath,
			want:           filepath.FromSlash(outputBasePath),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.metadata.Name, func(t *testing.T) {
			result, err := joinOutputBasePathWithRelativeInputPath(tc.inputBasePath, tc.inputFullPath, tc.outputBasePath)
			utils.TestingAssertErrorToWantError(t, err, tc.metadata.WantErr)
			if err == nil {
				utils.TestingAssertEqualStrings(t, tc.want, result)
			}
		})
	}
}
