package main

import (
	"strings"
	"testing"

	"github.com/markuskreukniet/markus-tools/go/utils/test"
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
	testCases := []test.TestCaseInput{
		test.TestingCreateTestCaseInput("Basic", input, false),
		test.TestingCreateTestCaseInput("Empty Input", "", false),
	}

	// run testCases
	for _, tc := range testCases {
		t.Run(tc.Metadata.Name, func(t *testing.T) {
			// arrange and teardown
			directories, fileSystemNodes := test.TestingCreateFilesAndDirectoriesByMultipleInputs(t, tc.Input)
			defer test.TestingRemoveDirectoryTrees(t, directories)
			var builder strings.Builder
			if len(directories) > 0 {
				isFirstWrite := true
				delimitedCommaStrings := test.TestingTrimSpaceTrimSuffixSplitOnSemicolonAndSort(tc.Input)
				for _, delimitedCommaString := range delimitedCommaStrings {
					directoryWithOptionalFileAsStrings := test.TestingTrimSpaceAndSplitOnComma(delimitedCommaString)
					if directoryWithOptionalFileAsStrings[3] != "" {

						// probably not optimal but results in less code, which is fine for testing
						if isFirstWrite {
							isFirstWrite = false
						} else {
							test.TestingWriteString(t, "\n\n", &builder)
						}

						test.TestingWriteString(t, directoryWithOptionalFileAsStrings[2]+"\n"+directoryWithOptionalFileAsStrings[3], &builder)
					}
				}
			}

			// act
			outcome, err := plainTextFilesToText(fileSystemNodes)

			// assert
			test.TestingAssertErrorToWantErrorAndOutcomeToBuilderString(t, err, tc.Metadata.WantErr, builder, outcome)
		})
	}
}
