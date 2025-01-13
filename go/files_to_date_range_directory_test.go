package main

import (
	"testing"

	"github.com/markuskreukniet/markus-tools/go/utils"
)

func createTestContent(subContent string) string {
	return "content" + subContent + "\ncontent" + subContent
}

func TestFilesToDateRangeDirectory(t *testing.T) {
	// test:
	//   removing empty directory in destination
	//   renaming directory in destination from a name with - to a name without -, and the other way around. With adding/removing a file and without
	//   removing a duplicate file that a directory in the destination directory already contains, and that a different directory already contains
	//   adding one and two files from source to destination in every sub directory that is a duplicate so instead remove the source file
	// V add a file from source to destination
	// V create a new destination directory

	// TODO: should also remove 0 byte files

	//   destination directory should also checks also its sub directory tree

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
	contentM100 := createTestContent("m 10 0")
	contentM101 := createTestContent("m 10 1")
	content10MN12 := createTestContent("10 M N 1 2") // TODO: only good naming
	// moveDuplicateFilesToNewDirectories := `
	// 	,2020-10-10T20:40:40Z,txt d n 0.txt,` + contentM100 + `;

	// 	directory 1,2020-10-20T20:40:40Z,txt d n 1.txt,` + contentM101 + `;
	// 	directory 1/directory 2,2020-10-20T20:40:40Z,txt d n 1 2.txt,` + contentM101 + `;
	// `

	// V trying to add duplicate files to existing directories
	contentME110 := createTestContent("m e 11 0")
	contentME111 := createTestContent("m e 11 1")
	contentME1112 := createTestContent("m e 11 1 2")
	content11 := createTestContent("11")
	content111 := createTestContent("11 1")
	content112 := createTestContent("11 2")
	content1112 := createTestContent("11 1 2")
	content113 := createTestContent("11 3")
	content1113 := createTestContent("11 1 3")
	content114 := createTestContent("11 4")
	content1114 := createTestContent("11 1 4")
	content115 := createTestContent("11 5")
	content116 := createTestContent("11 6")
	// moveDuplicateFilesToExistingDirectories := `
	// 	,2020-11-01T20:40:40Z,txt d e 0.txt,` + content11 + `;

	// 	directory 1,2020-11-06T20:40:40Z,txt d e 1.txt,` + content112 + `;
	// 	directory 1/directory 2,2020-11-06T20:40:40Z,txt d e 1 2.txt,` + content112 + `;
	// `

	input := ""
	destinationInput := ""
	wantedOutcome := ""

	// V move files to new directories
	// Moving files to a date range and a non-date range directory improves testing.
	input = input + `
		,2020-10-10T20:40:40Z,txt m n 0.txt,` + contentM100 + `;
		directory 1,2020-10-20T20:40:40Z,txt m n 1.txt,` + contentM101 + `;
		directory 1/directory 2,2020-10-22T20:40:40Z,txt m n 1 2.txt,` + content10MN12 + `;
	`
	wantedOutcome = wantedOutcome + `
		2020-10-10,2020-10-10T20:40:40Z,txt m n 0.txt,` + contentM100 + `;
		2020-10-20 - 2020-10-22,2020-10-20T20:40:40Z,txt m n 1.txt,` + contentM101 + `;
		2020-10-20 - 2020-10-22,2020-10-22T20:40:40Z,txt m n 1 2.txt,` + content10MN12 + `;
	`

	// V move files to existing directories
	input = input + `
		,2020-11-10T20:40:40Z,txt m e 0.txt,` + contentME110 + `;
		directory 1,2020-11-20T20:40:40Z,txt m e 1.txt,` + contentME111 + `;
		directory 1/directory 2,2020-11-21T20:40:40Z,txt m e 1 2.txt,` + contentME1112 + `;
	`
	wantedOutcome = wantedOutcome + `
		2020-11-06 - 2020-11-12,2020-11-10T20:40:40Z,txt m e 0.txt,` + contentME110 + `;
		2020-11-20 - 2020-11-23,2020-11-20T20:40:40Z,txt m e 1.txt,` + contentME111 + `;
		2020-11-20 - 2020-11-23,2020-11-21T20:40:40Z,txt m e 1 2.txt,` + contentME1112 + `;
	`

	// V removing empty directories and empty directory trees
	// V moving files to its parent directory
	// V renaming a date directory to a date range directory and a date range directory to a date directory
	// V having a date range directory to increase the date range of that directory
	destinationInput = destinationInput + `
		empty,,,;
		directory 1/empty,,,;
		2020-11-01,2020-11-01T20:40:40Z,txt 11.txt,` + content11 + `;
		2020-11-01/directory 1,2020-11-01T20:40:40Z,txt 11 1.txt,` + content111 + `;
		2020-11-01/empty,,,;
		2020-11-06 - 2020-11-07,2020-11-06T20:40:40Z,txt 11.txt,` + content112 + `;
		2020-11-06 - 2020-11-07/directory 1,2020-11-07T20:40:40Z,txt 11 1.txt,` + content1112 + `;
		2020-11-11,2020-11-11T20:40:40Z,txt 11.txt,` + content113 + `;
		2020-11-11/directory 1,2020-11-12T20:40:40Z,txt 11 1.txt,` + content1113 + `;
		2020-11-16 - 2020-11-17,2020-11-16T20:40:40Z,txt 11.txt,` + content114 + `;
		2020-11-16 - 2020-11-17/directory 1/directory 2,2020-11-16T20:40:40Z,txt 11 1.txt,` + content1114 + `;
		2020-11-21 - 2020-11-22,2020-11-21T20:40:40Z,txt 11.txt,` + content115 + `;
		2020-11-21 - 2020-11-22,2020-11-23T20:40:40Z,txt 11 2.txt,` + content116 + `;
	`
	wantedOutcome = wantedOutcome + `
		2020-11-01,2020-11-01T20:40:40Z,txt 11.txt,` + content11 + `;
		2020-11-01,2020-11-01T20:40:40Z,txt 11 1.txt,` + content111 + `;
		2020-11-06 - 2020-11-12,2020-11-06T20:40:40Z,txt 11.txt,` + content112 + `;
		2020-11-06 - 2020-11-12,2020-11-07T20:40:40Z,txt 11 1.txt,` + content1112 + `;
		2020-11-06 - 2020-11-12,2020-11-11T20:40:40Z,txt 11 2.txt,` + content113 + `;
		2020-11-06 - 2020-11-12,2020-11-12T20:40:40Z,txt 11 1 2.txt,` + content1113 + `;
		2020-11-16,2020-11-16T20:40:40Z,txt 11.txt,` + content114 + `;
		2020-11-16,2020-11-16T20:40:40Z,txt 11 1.txt,` + content1114 + `;
		2020-11-20 - 2020-11-23,2020-11-21T20:40:40Z,txt 11.txt,` + content115 + `;
		2020-11-20 - 2020-11-23,2020-11-23T20:40:40Z,txt 11 2.txt,` + content116 + `;
	`

	// V removing duplicate files in destination
	// V A three directories deep file (2020-12-21/directory 1/directory 2) improves testing
	content12 := createTestContent("12")
	content122 := createTestContent("12 2")
	destinationInput = destinationInput + `
		2020-12-20,2020-12-20T20:40:40Z,txt 12.txt,` + content12 + `;
		2020-12-20/directory 1,2020-12-20T20:40:40Z,txt 12 1.txt,` + content12 + `;
		directory 1,2020-12-21T20:40:40Z,txt 1.txt,` + content122 + `;
		2020-12-21,2020-12-21T20:40:40Z,txt 12 2.txt,` + content122 + `;
		2020-12-21/directory 1/directory 2,2020-12-22T20:40:40Z,txt 12 1 2.txt,` + content122 + `;
		2020-12-22,2020-12-22T20:40:40Z,txt 12 3.txt,` + content122 + `;
	`
	wantedOutcome = wantedOutcome + `
		2020-12-20 - 2020-12-21,,txt 1.txt,` + content122 + `;
		2020-12-20 - 2020-12-21,,txt 12.txt,` + content12 + `;
	`

	// TODO: "Empty Input" is missing
	testCases := []utils.TestCaseBasicDoubleInput{
		utils.CreateTestCaseBasicDoubleInput(utils.CreateTestCaseBasic("Basic", input, wantedOutcome, false), destinationInput),
	}

	// run testCases
	for _, tc := range testCases {
		t.Run(tc.TestCaseBasic.Name, func(t *testing.T) {
			// arrange and tear down
			directories, nodes := utils.WriteFilesByMultipleInputs(t, tc.TestCaseBasic.Input)
			defer utils.RemoveDirectoryTrees(t, directories)

			destination := utils.WriteFilesBySingleInput(t, tc.SecondInput)
			defer utils.TMustRemoveAll(t, destination)

			wantedOutcomeDestination := utils.WriteFilesBySingleInput(t, tc.TestCaseBasic.WantedOutcome)
			defer utils.TMustRemoveAll(t, wantedOutcomeDestination)

			// act
			err := filesToDateRangeDirectory(nodes, destination)

			// assert
			utils.TMustAssertError(t, err, tc.TestCaseBasic.WantErr)
			utils.TMustAssertIdenticalDescendantsFileTrees(t, destination, wantedOutcomeDestination)
		})
	}
}
