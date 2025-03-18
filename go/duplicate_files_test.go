package main

import (
	"path/filepath"
	"strings"
	"testing"

	"github.com/markuskreukniet/markus-tools/go/utils"
)

func TestGetDuplicateFilesAsNewlineSeparatedString(t *testing.T) {
	// arrange
	// Two contents should have the same file size and one a different file size.
	contents := []string{
		"content 1\ncontent 1",
		"content 2\ncontent 2",
		"content 3 1\ncontent 3 1",
	}
	input := `
		empty,,,;
		directory 2/empty,,,;
		directory 1,,txt 1.txt,;
		directory 1,,txt 1 2.txt,` + contents[0] + `;
		directory 2/directory 3,,txt 2-3.txt,` + contents[0] + `;
		directory 2/directory 3,,txt 2-3 2.txt,` + contents[1] + `;
		directory 2/directory 3,,txt 2-3 3.txt,` + contents[1] + `;
		directory 2/directory 4,,txt 2-4.txt,` + contents[1] + `;
		directory 5/directory 6/directory 7,,txt 5-6-7.txt,` + contents[2] + `;
		directory 8,,txt 8.txt,` + contents[2] + `;
	`
	wantedOutcome := `
		directory 1/txt 1 2.txt
		directory 2/directory 3/txt 2-3.txt

		directory 2/directory 3/txt 2-3 2.txt
		directory 2/directory 3/txt 2-3 3.txt
		directory 2/directory 4/txt 2-4.txt

		directory 5/directory 6/directory 7/txt 5-6-7.txt
		directory 8/txt 8.txt
	`
	testCases := []utils.TestCaseBasic{
		utils.CreateTestCaseBasic("Basic", input, wantedOutcome, false),
		utils.CreateTestCaseBasic("Empty Input", "", "", false),
	}

	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			// arrange and tear down
			directories, fileSystemNodes := utils.WriteFilesByMultipleInputs(t, tc.Input)
			defer utils.RemoveDirectoryTrees(t, directories)

			var trimmedLines []string
			for _, line := range strings.Split(wantedOutcome, "\n") {
				trimmed := strings.TrimSpace(line)
				if trimmed != "" {
					// 'filepath.FromSlash'
					// converts a slash-separated path to the native path format for the current operating system.
					trimmedLines = append(trimmedLines, filepath.FromSlash(trimmed))
				}
			}
			wantedOutcome = strings.Join(trimmedLines, "\n")

			// act
			outcome, err := getDuplicateFilesAsNewlineSeparatedString(fileSystemNodes)
			outcome = utils.TMust(t, outcome, err)

			// assert
			for _, directory := range directories {
				outcome = strings.ReplaceAll(outcome, directory+utils.FilePathSeparator, "")
			}

			for _, substring := range strings.Split(outcome, "\n\n") {
				wantedOutcome = strings.Replace(wantedOutcome, substring, "", 1)
			}

			utils.TMustAssertEqualStrings(t, strings.ReplaceAll(wantedOutcome, "\n\n", ""), "")
		})
	}
}
