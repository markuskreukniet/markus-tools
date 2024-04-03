package main

import (
	"strings"
	"testing"

	"github.com/markuskreukniet/markus-tools/go/utils/test"
)

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

	// run testCases
	for _, tc := range testCases {
		t.Run(tc.metadata.Name, func(t *testing.T) {
			// arrange and teardown
			directory, fileSystemNodes := test.TestingCreateFilesAndDirectories(t, tc.input)
			defer test.TestingRemoveDirectoryTree(t, directory)
			var builder strings.Builder
			if directory != "" {
				isFirstWrite := true
				for _, delimitedCommaString := range test.TestingTrimSpaceTrimSuffixOnSemicolonAndSplitOnSemicolon(tc.input) {
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
			test.TestingAssertErrorToWantErrorAndOutcomeToBuilderString(t, err, tc.metadata.WantErr, builder, outcome)
		})
	}
}
