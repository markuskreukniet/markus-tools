package main

import (
	"path/filepath"
	"strings"
	"testing"

	"github.com/markuskreukniet/markus-tools/go/utils"
	"github.com/markuskreukniet/markus-tools/go/utils/test"
)

// TODO: cleaning
func TestGetDuplicateFilesAsNewlineSeparatedString(t *testing.T) {
	// arrange
	// Two content strings should result in the same file size, and there should be a content string resulting in a different file size.
	contents := []string{
		"content 1\ncontent 1",
		"content 2\ncontent 2",
	}
	input := `
		empty,,,;
		directory 1,,txt 1.txt,;
		directory 1,,txt 1 2.txt,` + contents[0] + `;
		directory 2/empty,,,;
		directory 2/directory 3,,txt 2-3.txt,` + contents[0] + `;
		directory 2/directory 3,,txt 2-3 2.txt,` + contents[1] + `;
		directory 2/directory 3,,txt 2-3 3.txt,` + contents[1] + `;
		directory 2/directory 4,,txt 2-4.txt,` + contents[1] + `;
	`
	testCases := []test.TestCaseInput{
		test.TestingCreateTestCaseInput("Basic", input, false),
		test.TestingCreateTestCaseInput("Empty Input", "", false),
	}

	for _, tc := range testCases {
		t.Run(tc.Metadata.Name, func(t *testing.T) {
			// arrange and teardown
			directories, fileSystemNodes := test.TestingCreateFilesAndDirectoriesByMultipleInputs(t, tc.Input)
			defer test.TestingRemoveDirectoryTrees(t, directories)
			var builder strings.Builder

			// create duplicate file groups
			var fileGroups []duplicateFileGroup
			if len(directories) > 0 {
				var directoriesWithFileAsStrings [][]string
				for _, delimitedCommaString := range test.TestingTrimSpaceTrimSuffixSplitOnSemicolonAndSort(tc.Input) {
					inputLine := test.CreateInputLine(delimitedCommaString)
					if inputLine.HasNoContent() {
						directoriesWithFileAsStrings = append(directoriesWithFileAsStrings, inputLine)
						filePathPart := filepath.Join(inputLine[0], inputLine[2])
						for _, nodeI := range fileSystemNodes {
							if strings.HasSuffix(nodeI.Path, filePathPart) {
								// TODO: duplicate from non-test
								foundGroup := false
								for i, group := range fileGroups {
									if inputLine.GetContent() == group.identifier {
										foundGroup = true
										fileGroups[i].filePaths = append(fileGroups[i].filePaths, nodeI.Path)
										break
									}
								}
								if !foundGroup {
									for _, directoryWithFileAsStrings := range directoriesWithFileAsStrings {
										if inputLine.GetContent() == directoryWithFileAsStrings[3] {
											for _, nodeJ := range fileSystemNodes {
												if nodeI.Path != nodeJ.Path && strings.HasSuffix(nodeJ.Path, filepath.Join(directoryWithFileAsStrings[0], directoryWithFileAsStrings[2])) {
													fileGroups = append(fileGroups, duplicateFileGroup{
														identifier: directoryWithFileAsStrings[3],
														filePaths:  []string{nodeJ.Path, nodeI.Path},
													})
													break
												}
											}
											break
										}
									}
								}
								break
							}
						}

					}
				}

				// create and return the result string
				for i, group := range fileGroups {
					if i != 0 {
						test.TestingWriteTwoNewlineStrings(t, &builder)
					}
					for j, path := range group.filePaths {
						if j != 0 {
							if _, err := utils.WriteNewlineString(&builder); err != nil {
								t.Errorf("WriteNewlineString error: %v", err)
							}
						}
						test.TestingWriteString(t, path, &builder)
					}
				}
			}

			// act
			outcome, err := getDuplicateFilesAsNewlineSeparatedString(fileSystemNodes)

			// assert
			test.TestingAssertErrorToWantErrorAndOutcomeToBuilderString(t, err, tc.Metadata.WantErr, builder, outcome)
		})
	}
}
