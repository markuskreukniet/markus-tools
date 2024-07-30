package main

import (
	"strings"
	"testing"

	"github.com/markuskreukniet/markus-tools/go/utils"
)

// TODO: fix and clean
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
	wantedOutcome := "txt 0.txt\ncontent 0\\ncontent 0\n\ntxt 1-2 2.txt\ncontent directory 1/directory\\ncontent 2 1-2 2\n\ntxt 1-2.txt\ncontent directory 1/directory\\ncontent 2 1-2"
	testCases := []utils.TestCaseInput{
		utils.CreateTestCaseInput("Basic", input, false),
		utils.CreateTestCaseInput("Empty Input", "", false),
	}

	// run testCases
	for _, tc := range testCases {
		t.Run(tc.Metadata.Name, func(t *testing.T) {
			// arrange and teardown
			directories, fileSystemNodes := utils.TestingCreateFilesAndDirectoriesByMultipleInputs(t, tc.Input)
			defer utils.TestingRemoveDirectoryTrees(t, directories)
			var builder strings.Builder
			if len(directories) > 0 {
				builder.WriteString(wantedOutcome)
			}

			// act
			outcome, err := plainTextFilesToText(fileSystemNodes)

			//log.Println("outcome:", outcome) // TODO: shows a \n bug

			// assert
			utils.TestingAssertErrorToWantErrorAndOutcomeToBuilderString(t, err, tc.Metadata.WantErr, builder, outcome)
		})
	}
}
