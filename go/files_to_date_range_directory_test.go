package main

import (
	"os"
	"path/filepath"
	"sort"
	"strings"
	"testing"
	"time"

	"github.com/markuskreukniet/markus-tools/go/utils"
)

func createTestContent(subContent string) string {
	return "content" + subContent + "\ncontent" + subContent
}

func filterTimeInputLines(filteredLines *[]timeInputLine, filter func([]timeInputLine) []timeInputLine) {
	if len(*filteredLines) > 1 {
		tempLines := filter(*filteredLines)
		if len(tempLines) > 0 {
			*filteredLines = tempLines
		}
	}
}

type timeInputLine struct {
	time      time.Time
	inputLine utils.InputLine
}

func testingCreateTimeInputLine(t *testing.T, line utils.InputLine) timeInputLine {
	return timeInputLine{
		time:      utils.TestingParseTime(t, line.GetTimeModified()),
		inputLine: line,
	}
}

type duplicateTimeInputLineGroup struct {
	identifier     string
	timeInputLines []timeInputLine
}

func createDuplicateTimeInputLineGroup(identifier string, lines []timeInputLine) duplicateTimeInputLineGroup {
	return duplicateTimeInputLineGroup{
		identifier:     identifier,
		timeInputLines: lines,
	}
}

type duplicateTimeInputLineGroups []duplicateTimeInputLineGroup

func (groups duplicateTimeInputLineGroups) appendByIdentifier(identifier string, line timeInputLine) bool {
	for i, group := range groups {
		if identifier == group.identifier {
			groups[i].timeInputLines = append(groups[i].timeInputLines, line)
			return true
		}
	}
	return false
}

// garbage collection: unGroupedLines
func createDuplicateInputLineFileGroups(t *testing.T, input string) duplicateTimeInputLineGroups {
	var groups duplicateTimeInputLineGroups
	var unGroupedLines []timeInputLine
	for _, rawLine := range utils.CreateSortedRawInputLines(input) {
		line := utils.CreateInputLine(rawLine)

		if !line.HasContent() {
			continue
		}

		timeLine := testingCreateTimeInputLine(t, line)
		if groups.appendByIdentifier(line.GetContent(), timeLine) {
			continue
		}

		var lines []timeInputLine
		for _, unGroupedLine := range unGroupedLines {
			if line.GetContent() == unGroupedLine.inputLine.GetContent() {
				lines = append(lines, unGroupedLine)
			}
		}
		if len(lines) > 0 {
			groups = append(groups, createDuplicateTimeInputLineGroup(line.GetContent(), append([]timeInputLine{timeLine}, lines...)))
		} else {
			unGroupedLines = append(unGroupedLines, timeLine)
		}
	}
	return groups
}

func directoryTreeToFilePathHashes(directory string) (map[string]string, error) {
	hashes := make(map[string]string)

	err := filepath.Walk(directory, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			hash, err := utils.HashFile(path)
			if err != nil {
				return err
			}
			hashes[path] = hash
		}
		return nil
	})

	return hashes, err
}

func areDirectoryTreesTheSame(dir1, dir2 string) (bool, error) {
	hashes1, err := directoryTreeToFilePathHashes(dir1)
	if err != nil {
		return false, err
	}

	hashes2, err := directoryTreeToFilePathHashes(dir2)
	if err != nil {
		return false, err
	}

	for path1, hash1 := range hashes1 {
		relativePath, err := filepath.Rel(dir1, path1)
		if err != nil {
			return false, err
		}
		path2 := filepath.Join(dir2, relativePath)
		hash2, found := hashes2[path2]
		if !found || hash1 != hash2 {
			return false, nil
		}
		delete(hashes2, path2)
	}

	if len(hashes2) > 0 {
		return false, nil
	}

	return true, nil
}

//

func TestFilesToDateRangeDirectory(t *testing.T) {
	// test:
	//   removing empty folder in destination
	//   renaming folder in destination from a name with - to a name without -, and the other way around. With adding/removing a file and without
	//   removing a duplicate file that a folder in the destination folder already contains, and that a different folder already contains
	//   adding one and two files from source to destination in every sub folder that is a duplicate so instead remove the source file
	// V add a file from source to destination
	// V create a new destination folder

	// TODO: should also remove 0 byte files

	//   destination folder should also checks also its sub directory tree

	// add a file with the same file name, but first rename it with 2 postfix. Or replace file if newer
	// Maximaal 255 tekens voor een bestandsnaam.
	// De volledige padlengte (inclusief de mappenstructuur) is beperkt tot 260 tekens.

	// arrange
	// time022040 := "2020-02-20T20:40:40Z"
	// time022041 := "2020-02-20T20:40:41Z"
	// time0520 := "2020-05-20T20:40:40Z"
	// time0521 := "2020-05-21T20:40:40Z"

	// // // //

	// // dirs wrong date name named test

	// // for changing (creating a new directory) a date directory to a date range directory
	// time0104 := "2020-01-21T20:40:40Z"

	// // for increasing (creating a new directory) a the date range of a directory
	// time0224 := "2020-02-24T20:40:40Z"

	// // for changing (creating a new directory) a date range directory to a date directory
	// time0321 := "2020-03-21T20:40:40Z"

	// // for decreasing (creating a new directory) a the date range of a directory, can't?
	// time0425 := "2020-04-25T20:40:40Z"

	// // // //

	// // moving files to rename a new date directories to date range directories
	// // moving files to rename an existing date directories to date range directories

	// // moving files to rename a new date range directories to a bigger range

	// // V moving files to rename an existing date range directories to a bigger range
	// moveFilesToIncreaseDateRangeOfExistingDirectories := `
	// 	,2020-11-24T20:40:40Z,txt m e 0.txt,` + TestingCreateContent("m e 0") + `;
	// 	directory 1,2020-11-24T20:40:40Z,txt m e 1.txt,` + TestingCreateContent("m e 1") + `;
	// `

	// // V trying to add a duplicate files to new directories
	// contentM100 := TestingCreateContent("m 10 0")
	// contentM101 := TestingCreateContent("m 10 1")
	// moveDuplicateFilesToNewDirectories := `
	// 	,2020-10-10T20:40:40Z,txt d n 0.txt,` + contentM100 + `;

	// 	directory 1,2020-10-20T20:40:40Z,txt d n 1.txt,` + contentM101 + `;
	// 	directory 1/directory 2,2020-10-20T20:40:40Z,txt d n 1 2.txt,` + contentM101 + `;
	// `

	// // V trying to add duplicate files to existing directories
	// content11 := TestingCreateContent("11")
	// content112 := TestingCreateContent("11 2")
	// moveDuplicateFilesToExistingDirectories := `
	// 	,2020-11-01T20:40:40Z,txt d e 0.txt,` + content11 + `;

	// 	directory 1,2020-11-06T20:40:40Z,txt d e 1.txt,` + content112 + `;
	// 	directory 1/directory 2,2020-11-06T20:40:40Z,txt d e 1 2.txt,` + content112 + `;
	// `

	// // V move files to new directories
	// // --- move files to a new date range directory to increase the date range of that directory with a different variable
	// moveFilesToNewDirectories := `
	// 	,2020-10-10T20:40:40Z,txt m n 0.txt,` + contentM100 + `;

	// 	directory 1,2020-10-20T20:40:40Z,txt m n 1.txt,` + contentM101 + `;
	// 	directory 1/directory 2,2020-10-22T20:40:40Z,txt m n 1 2.txt,` + TestingCreateContent("m n 10 1 2") + `;
	// `

	// // V move files to existing directories
	// moveFilesToExistingDirectories := `
	// 	,2020-11-10T20:40:40Z,txt m e 0.txt,` + TestingCreateContent("m e 11 0") + `;

	// 	directory 1,2020-11-20T20:40:40Z,txt m e 1.txt,` + TestingCreateContent("m e 11 1") + `;
	// 	directory 1/directory 2,2020-11-21T20:40:40Z,txt m e 1 2.txt,` + TestingCreateContent("m e 11 1 2") + `;
	// `

	// // V removing empty directories and empty directory trees
	// // V moving files to its parent directory
	// // V renaming a date directory to a date range directory and a date range directory to a date directory
	// // V having a date range directory to increase the date range of that directory with a different variable
	// fixingDestinationBesidesDuplicates := `
	// 	empty,,,;
	// 	directory 1/empty,,,;

	// 	2020-11-01,2020-11-01T20:40:40Z,txt 11.txt,` + content11 + `;
	// 	2020-11-01/directory 1,2020-11-01T20:40:40Z,txt 11 1.txt,` + TestingCreateContent("11 1") + `;
	// 	2020-11-01/empty,,,;

	// 	2020-11-06 - 2020-11-07,2020-11-06T20:40:40Z,txt 11.txt,` + content112 + `;
	// 	2020-11-06 - 2020-11-07/directory 1,2020-11-07T20:40:40Z,txt 11 1.txt,` + TestingCreateContent("11 1 2") + `;

	// 	2020-11-11,2020-11-11T20:40:40Z,txt 11.txt,` + TestingCreateContent("11 3") + `;
	// 	2020-11-11/directory 1,2020-11-12T20:40:40Z,txt 11 1.txt,` + TestingCreateContent("11 1 3") + `;

	// 	2020-11-16 - 2020-11-17,2020-11-16T20:40:40Z,txt 11.txt,` + TestingCreateContent("11 4") + `;
	// 	2020-11-16 - 2020-11-17/directory 1,2020-11-16T20:40:40Z,txt 11 1.txt,` + TestingCreateContent("11 1 4") + `;

	// 	2020-11-21 - 2020-11-22,2020-11-21T20:40:40Z,txt 11.txt,` + TestingCreateContent("11 5") + `;
	// 	2020-11-21 - 2020-11-22,2020-11-22T20:40:40Z,txt 11.txt,` + TestingCreateContent("11 6") + `;
	// `

	// V removing duplicate files in destination
	// V A three directories deep file improves testing
	content12 := createTestContent("12")
	content122 := createTestContent("12 2")
	destinationDuplicateFiles := `
		2020-12-20,2020-12-20T20:40:40Z,txt 12.txt,` + content12 + `;
		2020-12-20/directory 1,2020-12-20T20:40:40Z,txt 12 1.txt,` + content12 + `;

		directory 1,2020-12-21T20:40:40Z,txt 1.txt,` + content122 + `;
		2020-12-21,2020-12-21T20:40:40Z,txt 12 2.txt,` + content122 + `;
		2020-12-21/directory 1/directory 2,2020-12-22T20:40:40Z,txt 12 1 2.txt,` + content122 + `;
		2020-12-22,2020-12-22T20:40:40Z,txt 12 3.txt,` + content122 + `;
	`

	input := ""
	destinationInput := destinationDuplicateFiles

	// create testCases
	testCases := []struct {
		testCaseInput    utils.TestCaseInput
		destinationInput string
	}{
		{
			testCaseInput:    utils.CreateTestCaseInput("Basic", input, false),
			destinationInput: destinationInput,
		},
		{
			testCaseInput:    utils.CreateTestCaseInput("Empty Input", "", false),
			destinationInput: destinationInput,
		},
	}

	// run testCases
	for _, tc := range testCases {
		t.Run(tc.testCaseInput.Metadata.Name, func(t *testing.T) {
			// arrange and teardown
			directory, _ := utils.TestingCreateFilesAndDirectoriesByOneInput(t, tc.destinationInput)
			defer utils.TestingRemoveDirectoryTree(t, directory)

			destination := utils.CreateTemporaryDirectory(t)
			defer utils.TestingRemoveDirectoryTree(t, destination)

			groups := createDuplicateInputLineFileGroups(t, tc.destinationInput)

			// Select unique files by first filtering the duplicate ones (there are no created files yet) by this priority:
			// 1. keep the shortest file name
			// 2. keep the one in the destination a date directory or date range directory
			// 3. keep the one in the destination directory
			// 4. keep the newest modification time file
			// 5. keep the first file of the slice
			var lines []timeInputLine
			for _, group := range groups {
				// TODO: It is possible to clean the anonymous functions in filterTimeInputLines
				if len(group.timeInputLines) > 1 {
					// filter on shortest file name
					filterTimeInputLines(&group.timeInputLines, func(unfilteredLines []timeInputLine) []timeInputLine {
						var tempLines []timeInputLine
						var minimumLength int
						for _, line := range unfilteredLines {
							length := len(line.inputLine.GetFileName())
							if length < minimumLength || minimumLength == 0 {
								minimumLength = length
								tempLines = []timeInputLine{line}
							} else if length == minimumLength {
								tempLines = append(tempLines, line)
							}
						}
						return tempLines
					})

					// filter on valid name of date directory or date range directory
					filterTimeInputLines(&group.timeInputLines, func(unfilteredLines []timeInputLine) []timeInputLine {
						var tempLines []timeInputLine
						for _, line := range unfilteredLines {
							part := line.inputLine.GetDirectoryPathPart()
							slash := "/"
							if strings.Contains(part, slash) {
								// TODO: is this correct?
								subStrings := strings.SplitN(part, slash, 2)
								if len(subStrings) > 0 {
									part = subStrings[0]
								}
							}
							if isValidDateRangeDirectoryName(part) {
								tempLines = append(tempLines, line)
							}
						}
						return tempLines
					})

					// filter on destination directory

					// filter on the newest modification time file
					filterTimeInputLines(&group.timeInputLines, func(unfilteredLines []timeInputLine) []timeInputLine {
						var tempLines []timeInputLine
						var newestTime time.Time
						for _, line := range unfilteredLines {
							if line.time.After(newestTime) {
								newestTime = line.time
								tempLines = []timeInputLine{line}
							} else if line.time.Equal(newestTime) {
								tempLines = append(tempLines, line)
							}
						}
						return tempLines
					})
				}

				// keep the first file of the slice
				lines = append(lines, group.timeInputLines[0])
			}

			// create unique files in directories
			sort.Slice(lines, func(i, j int) bool {
				return lines[i].time.Before(lines[j].time)
			})
			startDateRange := 0
			isFindingDateRange := false
			length := len(lines)
			for i := 0; i < length; i++ {
				if i < length-1 && isWithinThreeDays(lines[i].time, lines[i+1].time) && !isFindingDateRange {
					isFindingDateRange = true
					startDateRange = i
				} else {
					var name string
					if isFindingDateRange {
						// Declare 'err' separately to avoid shadowing 'name' with ':='
						var err error
						name, err = createDirectoryDateRangeName(lines[startDateRange].time, lines[i].time)
						if err != nil {
							t.Errorf("createDirectoryDateRangeName error: %v", err)
						}

						isFindingDateRange = false
					} else {
						name = toDateFormat(lines[i].time)
					}

					// create directory with files
					path := filepath.Join(destination, name)
					utils.TestingCreateDirectoryAll(t, path)
					for j := startDateRange; j <= i; j++ {
						// TODO: write also time modified?
						utils.TestingWriteFile(t, filepath.Join(path, lines[j].inputLine.GetFileName()), lines[j].inputLine.GetContent())
					}
				}
			}

			// act
			if err := filesToDateRangeDirectory(nil, directory); err != nil {
				t.Errorf("filesToDateRangeDirectory error: %v", err)
			}

			if same, err := areDirectoryTreesTheSame(directory, destination); err != nil {
				t.Errorf("compareDirectories error: %v", err)
				t.Errorf("compareDirectories same: %v", same)
			}

			// assert
		})
	}
}
