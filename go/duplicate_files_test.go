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

			if len(directories) > 0 {
				// create duplicate file groups
				var groups utils.DuplicateFileGroups
				var inputLines []utils.InputLine
				for _, rawInputLine := range utils.CreateSortedRawInputLines(tc.Input) {
					line := utils.CreateInputLine(rawInputLine)

					if !line.HasContent() {
						continue
					}

					for _, nodeI := range fileSystemNodes {
						if !strings.HasSuffix(nodeI.Path, line.JoinDirectoryPathPartWithFileName()) {
							continue
						}

						appended := groups.AppendByIdentifier(line.GetContent(), nodeI.Path)
						if !appended {
							var paths []string
							for _, unGroupedLine := range inputLines {
								if line.GetContent() == unGroupedLine.GetContent() {
									for _, nodeJ := range fileSystemNodes {
										// TODO: JoinDirectoryPathPartWithFileName should happen only when it is not done before
										if nodeI.Path != nodeJ.Path && strings.HasSuffix(nodeJ.Path, unGroupedLine.JoinDirectoryPathPartWithFileName()) {
											paths = append(paths, nodeJ.Path)
											break
										}
									}
								}
							}
							if len(paths) > 0 {
								groups = append(groups, utils.DuplicateFileGroup{
									Identifier: line.GetContent(),
									FilePaths:  append(paths, []string{nodeI.Path}...),
								})
							} else {
								inputLines = append(inputLines, line)
							}
						}
						break
					}
				}

				// create and return the result string
				for i, group := range groups {
					if i != 0 {
						utils.TestingWriteTwoNewlineStrings(t, &builder)
					}
					for j, path := range group.FilePaths {
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
