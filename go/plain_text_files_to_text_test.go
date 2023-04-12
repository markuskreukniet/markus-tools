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
	wantedOutcome := "txt 1-2 2.txt\ncontent directory 1/directory\\ncontent 2 1-2 2\n\ntxt 1-2.txt\ncontent directory 1/directory\\ncontent 2 1-2\n\ntxt 0.txt\ncontent 0\\ncontent 0"
	testCases := []utils.TestCaseBasic{
		utils.CreateTestCaseBasic("Basic", input, wantedOutcome, false),
		utils.CreateTestCaseBasic("Empty Input", "", "", false),
	}

	// run testCases
	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			// arrange and tear down
			directories, fileSystemNodes := utils.WriteFilesByMultipleInputs(t, tc.Input)
			defer utils.RemoveDirectoryTrees(t, directories)
			var builder strings.Builder
			if len(directories) > 0 {
				builder.WriteString(tc.WantedOutcome)
			}

			// act
			outcome, err := plainTextFilesToText(fileSystemNodes)

			//log.Println("outcome:", outcome) // TODO: shows a \n bug, but it is nog a bug?

			// assert
			utils.TMustAssertError(t, err, tc.WantErr)
			utils.TMustAssertEqualStrings(t, builder.String(), outcome)
		})
	}
}
