package main

import (
	"path/filepath"
	"strings"
	"testing"

	"github.com/markuskreukniet/markus-tools/go/utils"
	"github.com/markuskreukniet/markus-tools/go/utils/test"
)

// TODO: same as fileGroup or duplicateFileGroup from non-test
type testingDuplicateFileGroup struct {
	content   string
	filePaths []string
}

func testingWriteNewlineString(t *testing.T, builder *strings.Builder) {
	t.Helper()
	if _, err := utils.WriteNewlineString(builder); err != nil {
		t.Errorf("writeNewlineString error: %v", err)
	}
}

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
	testCases := []struct {
		metadata test.TestCaseMetadata
		input    string
		contents []string
	}{
		{
			metadata: test.TestingCreateTestCaseMetadataWithNameBasicAndWantErrFalse(),
			input:    input,
			contents: contents,
		},
		{
			metadata: test.TestingCreateTestCaseMetadataWithNameEmptyFileSystemNodesAndWantErrFalse(),
			input:    "",
			contents: make([]string, 0),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.metadata.Name, func(t *testing.T) {
			// arrange and teardown
			directories, fileSystemNodes := test.TestingCreateFilesAndDirectoriesByMultipleInputs(t, tc.input)
			defer test.TestingRemoveDirectoryTrees(t, directories)
			var builder strings.Builder

			// create duplicate file groups
			var fileGroups []testingDuplicateFileGroup
			// TODO: needed?
			if len(directories) > 0 {
				var directoriesWithFileAsStrings [][]string
				for _, delimitedCommaString := range test.TestingTrimSpaceTrimSuffixSplitOnSemicolonAndSort(tc.input) {
					directoryWithOptionalFileAsStrings := test.TestingTrimSpaceAndSplitOnComma(delimitedCommaString)
					if directoryWithOptionalFileAsStrings[3] != "" {
						directoriesWithFileAsStrings = append(directoriesWithFileAsStrings, directoryWithOptionalFileAsStrings)
						filePathPart := filepath.Join(directoryWithOptionalFileAsStrings[0], directoryWithOptionalFileAsStrings[2])

						for _, nodeI := range fileSystemNodes {
							if strings.HasSuffix(nodeI.Path, filePathPart) {

								// TODO: duplicate from non-test
								foundGroup := false
								for i, group := range fileGroups {
									if directoryWithOptionalFileAsStrings[3] == group.content {
										foundGroup = true
										fileGroups[i].filePaths = append(fileGroups[i].filePaths, nodeI.Path)
										break
									}
								}

								// check
								if !foundGroup {
									for _, directoryWithFileAsStrings := range directoriesWithFileAsStrings {
										if directoryWithOptionalFileAsStrings[3] == directoryWithFileAsStrings[3] {
											for _, nodeJ := range fileSystemNodes {
												// ?
												if nodeI.Path != nodeJ.Path && strings.HasSuffix(nodeJ.Path, filepath.Join(directoryWithFileAsStrings[0], directoryWithFileAsStrings[2])) {
													fileGroups = append(fileGroups, testingDuplicateFileGroup{
														content:   directoryWithFileAsStrings[3],
														filePaths: []string{nodeJ.Path, nodeI.Path},
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
				//

				// create and return the result string
				//var result strings.Builder
				for i, group := range fileGroups {
					if i != 0 {
						if _, err := utils.WriteTwoNewlineStrings(&builder); err != nil {
							t.Errorf("TODO")
						}
					}
					for j, path := range group.filePaths {
						if j != 0 {
							if _, err := utils.WriteNewlineString(&builder); err != nil {
								t.Errorf("TODO")
							}
						}
						if _, err := builder.WriteString(path); err != nil {
							t.Errorf("TODO")
						}
					}
				}
			}

			// act
			outcome, err := getDuplicateFilesAsNewlineSeparatedString(fileSystemNodes)

			// assert
			test.TestingAssertErrorToWantErrorAndOutcomeToBuilderString(t, err, tc.metadata.WantErr, builder, outcome)
		})
	}
}
