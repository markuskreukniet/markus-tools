package main

import (
	"path/filepath"
	"sort"
	"strings"
	"testing"
	"time"

	"github.com/markuskreukniet/markus-tools/go/utils"
)

// func getDirectoryPathsAndNamesFromDirectory(directoryFilePath string) ([]string, []string, error) {
// 	entries, err := os.ReadDir(directoryFilePath)
// 	if err != nil {
// 		return nil, nil, err
// 	}
// 	var filePaths []string
// 	var names []string
// 	for _, entry := range entries {
// 		if entry.IsDir() {
// 			names = append(names, entry.Name())
// 			filePaths = append(filePaths, filepath.Join(directoryFilePath, entry.Name()))
// 		}
// 	}
// 	return filePaths, names, nil
// }

func testingCreateContent(subContent string) string {
	return "content" + subContent + "\ncontent" + subContent
}

func testingCreateFileDetail(t *testing.T, line utils.InputLine, destination string) utils.FileDetail {
	time, err := time.Parse(time.RFC3339, line.GetTimeModified()) // TODO rename to ModificationTime
	if err != nil {
		t.Errorf("time.Parse error: %v", err)
	}

	return utils.FileDetail{
		Path:             filepath.Join(destination, line.GetDirectoryPathPartWithFileName()),
		ModificationTime: time,
		Size:             0,
	}
}

func testingFilterFileDetails(filteredDetails *[]utils.FileDetail, filter func(unfilteredDetails []utils.FileDetail) []utils.FileDetail) {
	if len(*filteredDetails) > 1 {
		tempDetails := filter(*filteredDetails)
		if len(tempDetails) > 1 {
			*filteredDetails = tempDetails
		}
	}
}

type duplicateFileDetailGroup struct {
	identifier  string
	fileDetails []utils.FileDetail
}

type duplicateFileDetailGroups []duplicateFileDetailGroup

func (groups duplicateFileDetailGroups) appendByIdentifier(identifier string, detail utils.FileDetail) bool {
	for i, group := range groups {
		if identifier == group.identifier {
			groups[i].fileDetails = append(groups[i].fileDetails, detail)
			return true
		}
	}
	return false
}

func TestFilesToDateRangeDirectory(t *testing.T) {
	// test:
	//   removing empty folder in destination
	//   renaming folder in destination from a name with - to a name without -, and the other way around. With adding/removing a file and without
	//   removing a duplicate file that a folder in the destination folder already contains, and that a different folder already contains
	//   adding one and two files from source to destination in every sub folder that is a duplicate so instead remove the source file
	// V add a file from source to destination
	// V create a new destination folder

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
	content12 := testingCreateContent("12")
	content122 := testingCreateContent("12 2")
	destinationDuplicateFiles := `
		2020-12-20,2020-12-20T20:40:40Z,txt 12.txt,` + content12 + `;
		2020-12-20/directory 1,2020-12-20T20:40:40Z,txt 12 1.txt,` + content12 + `;

		directory 1,2020-12-21T20:40:40Z,txt 1.txt,` + content122 + `;
		2020-12-21,2020-12-21T20:40:40Z,txt 12 2.txt,` + content122 + `;
		2020-12-21/directory 1,2020-12-22T20:40:40Z,txt 12 1 2.txt,` + content122 + `;
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
			testCaseInput:    utils.TestingCreateTestCaseInput("Basic", input, false),
			destinationInput: destinationInput,
		},
		{
			testCaseInput:    utils.TestingCreateTestCaseInput("Empty Input", "", false),
			destinationInput: destinationInput,
		},
	}

	// run testCases
	for _, tc := range testCases {
		t.Run(tc.testCaseInput.Metadata.Name, func(t *testing.T) {
			// arrange and teardown
			destination := utils.CreateTemporaryDirectory(t)
			defer utils.TestingRemoveDirectoryTree(t, destination)

			// create duplicate file groups
			var groups duplicateFileDetailGroups
			var unGroupedLines []utils.InputLine
			for _, rawLine := range utils.CreateSortedRawInputLines(tc.destinationInput) {
				line := utils.CreateInputLine(rawLine)

				if !line.HasContent() {
					continue
				}

				detail := testingCreateFileDetail(t, line, destination)
				appended := groups.appendByIdentifier(line.GetContent(), detail)

				if appended {
					continue
				}

				var details []utils.FileDetail
				for _, unGroupedLine := range unGroupedLines {
					if line.GetContent() == unGroupedLine.GetContent() {
						details = append(details, testingCreateFileDetail(t, unGroupedLine, destination))
					}
				}
				if len(details) > 0 {
					groups = append(groups, duplicateFileDetailGroup{
						identifier:  line.GetContent(),
						fileDetails: append([]utils.FileDetail{detail}, details...),
					})
				} else {
					unGroupedLines = append(unGroupedLines, line)
				}
			}

			// Select unique files by first filtering the duplicate ones (there are no created files yet) by this priority:
			// 1. keep the shortest file name
			// 2. keep the one in the destination a date directory or date range directory
			// 3. keep the one in the destination directory
			// 4. keep the newest modification time file
			// 5. keep the first file of the slice
			var details []utils.FileDetail
			for _, group := range groups {
				// TODO: It is possible to clean the anonymous functions in testingFilterFileDetails
				if len(group.fileDetails) > 1 {
					// filter on shortest file name
					testingFilterFileDetails(&group.fileDetails, func(unfilteredDetails []utils.FileDetail) []utils.FileDetail {
						var tempDetails []utils.FileDetail
						var minimumLength int
						for _, detail := range unfilteredDetails {
							length := len(filepath.Base(detail.Path)) // TODO: is this efficient?
							if length < minimumLength {
								minimumLength = length
								tempDetails = []utils.FileDetail{detail}
							} else if length == minimumLength {
								tempDetails = append(tempDetails, detail)
							}
						}
						return tempDetails
					})

					// filter on valid name of date directory or date range directory
					testingFilterFileDetails(&group.fileDetails, func(unfilteredDetails []utils.FileDetail) []utils.FileDetail {
						var tempDetails []utils.FileDetail
						for _, detail := range unfilteredDetails {
							if isValidDateRangeDirectory(detail.Path) {
								tempDetails = append(tempDetails, detail)
							}
						}
						return tempDetails
					})

					// filter on destination directory

					// filter on the newest modification time file
					testingFilterFileDetails(&group.fileDetails, func(unfilteredDetails []utils.FileDetail) []utils.FileDetail {
						var tempDetails []utils.FileDetail
						var newestTime time.Time
						for _, detail := range unfilteredDetails {
							time := detail.ModificationTime
							if time.After(newestTime) {
								newestTime = time
								tempDetails = []utils.FileDetail{detail}
							} else if time.Equal(newestTime) {
								tempDetails = append(tempDetails, detail)
							}
						}
						return tempDetails
					})
				}

				// keep the first file of the slice
				details = append(details, group.fileDetails[0])
			}

			// unique files to date range groups
			sort.Slice(details, func(i, j int) bool {
				return details[i].ModificationTime.Before(details[j].ModificationTime)
			})
			startDateRange := 0
			isFindingDateRange := false
			for i := 1; i < len(details); i++ {
				if isWithinThreeDays(details[i-1].ModificationTime, details[i].ModificationTime) {
					isFindingDateRange = true
				} else {
					var name string
					if isFindingDateRange {
						// Declare err separately to avoid shadowing with ':='
						var err error
						name, err = createDirectoryDateRangeName(details[startDateRange].ModificationTime, details[i].ModificationTime)
						if err != nil {
							t.Errorf("createDirectoryDateRangeName error: %v", err)
						}
						isFindingDateRange = false
						startDateRange = i + 1
					} else {
						name = toDateFormat(details[i].ModificationTime)
						startDateRange++
					}

					// create directory

					utils.TestingCreateDirectoryAll(t, filepath.Join(destination, name))

					// // add files to directory
					// for j := startDateRange; j <= i; j++ {

					// }
				}
			}

			// act
			// assert
		})
	}
}

// TODO: use also in non test version?
func createDirectoryDateRangeName(startTime, endTime time.Time) (string, error) {
	var builder strings.Builder
	if err := formatDateAndWriteString(&builder, startTime); err != nil {
		return "", err
	}
	if _, err := builder.WriteString(spacedHyphen); err != nil {
		return "", err
	}
	if err := formatDateAndWriteString(&builder, endTime); err != nil {
		return "", err
	}
	return builder.String(), nil
}
