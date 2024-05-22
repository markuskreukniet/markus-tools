package main

import (
	"strings"
	"testing"

	"github.com/markuskreukniet/markus-tools/go/utils"
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
	testCases := []utils.TestCaseInput{
		utils.TestingCreateTestCaseInput("Basic", input, false),
		utils.TestingCreateTestCaseInput("Empty Input", "", false),
	}

	for _, tc := range testCases {
		t.Run(tc.Metadata.Name, func(t *testing.T) {
			// arrange and teardown
			directories, fileSystemNodes := utils.TestingCreateFilesAndDirectoriesByMultipleInputs(t, tc.Input)
			defer utils.TestingRemoveDirectoryTrees(t, directories)
			var builder strings.Builder

			// create duplicate file groups
			var fileGroups []duplicateFileGroup
			if len(directories) > 0 {
				var inputLines []utils.InputLine
				for _, rawInputLine := range utils.CreateSortedRawInputLines(tc.Input) {
					inputLine := utils.CreateInputLine(rawInputLine)
					if inputLine.HasContent() {
						inputLines = append(inputLines, inputLine)
						for _, nodeI := range fileSystemNodes {
							if strings.HasSuffix(nodeI.Path, inputLine.JoinDirectoryPathPartWithFileName()) {
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
									for _, line := range inputLines {
										if inputLine.GetContent() == line.GetContent() {
											for _, nodeJ := range fileSystemNodes {
												if nodeI.Path != nodeJ.Path && strings.HasSuffix(nodeJ.Path, line.JoinDirectoryPathPartWithFileName()) {
													fileGroups = append(fileGroups, duplicateFileGroup{
														identifier: line.GetContent(),
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
						utils.TestingWriteTwoNewlineStrings(t, &builder)
					}
					for j, path := range group.filePaths {
						if j != 0 {
							if _, err := utils.WriteNewlineString(&builder); err != nil {
								t.Errorf("WriteNewlineString error: %v", err)
							}
						}
						utils.TestingWriteString(t, path, &builder)
					}
				}
			}

			// act
			outcome, err := getDuplicateFilesAsNewlineSeparatedString(fileSystemNodes)

			// assert
			utils.TestingAssertErrorToWantErrorAndOutcomeToBuilderString(t, err, tc.Metadata.WantErr, builder, outcome)
		})
	}
}
