package main

import (
	"strings"
	"testing"

	"github.com/markuskreukniet/markus-tools/go/utils"
)

func TestPlainTextFilesToText(t *testing.T) {
	// arrange
	input := `
		,,txt 0.txt,content 0\ncontent 0;
		,,jpg 0.jpg,;
		empty,,,;
		directory 1/empty,,,;
		directory 1/directory 2,,txt 1-2.txt,content directory 1/directory\ncontent 2 1-2;
		directory 1/directory 2,,txt 1-2 2.txt,content directory 1/directory\ncontent 2 1-2 2;
	`
	testCases := []utils.TestCaseInput{
		utils.TestingCreateTestCaseInput("Basic", input, false),
		utils.TestingCreateTestCaseInput("Empty Input", "", false),
	}

	// run testCases
	for _, tc := range testCases {
		t.Run(tc.Metadata.Name, func(t *testing.T) {
			// arrange and teardown
			directories, fileSystemNodes := utils.TestingCreateFilesAndDirectoriesByMultipleInputs(t, tc.Input)
			defer utils.TestingRemoveDirectoryTrees(t, directories)
			var builder strings.Builder
			if len(directories) > 0 {
				isFirstWrite := true
				delimitedCommaStrings := utils.TestingTrimSpaceTrimSuffixSplitOnSemicolonAndSort(tc.Input)
				for _, delimitedCommaString := range delimitedCommaStrings {
					inputLine := utils.CreateInputLine(delimitedCommaString)
					if inputLine.HasNoContent() {

						// probably not optimal but results in less code, which is fine for testing
						if isFirstWrite {
							isFirstWrite = false
						} else {
							utils.TestingWriteTwoNewlineStrings(t, &builder)
						}

						utils.TestingWriteString(t, inputLine.GetFileName()+"\n"+inputLine.GetContent(), &builder)
					}
				}
			}

			// act
			outcome, err := plainTextFilesToText(fileSystemNodes)

			// assert
			utils.TestingAssertErrorToWantErrorAndOutcomeToBuilderString(t, err, tc.Metadata.WantErr, builder, outcome)
		})
	}
}
