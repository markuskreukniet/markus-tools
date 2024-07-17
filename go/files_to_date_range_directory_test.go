package main

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/markuskreukniet/markus-tools/go/utils"
)

func createTestContent(subContent string) string {
	return "content" + subContent + "\ncontent" + subContent
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

// TODO: does not checks for similar modification times?
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

	// V trying to add duplicate files to existing directories
	// contentME110 := createTestContent("m e 11 0")
	// contentME111 := createTestContent("m e 11 1")
	// contentME1112 := createTestContent("m e 11 1 2")
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

	// // V move files to new directories
	// // --- move files to a new date range directory to increase the date range of that directory with a different variable
	// moveFilesToNewDirectories := `
	// 	,2020-10-10T20:40:40Z,txt m n 0.txt,` + contentM100 + `;

	// 	directory 1,2020-10-20T20:40:40Z,txt m n 1.txt,` + contentM101 + `;
	// 	directory 1/directory 2,2020-10-22T20:40:40Z,txt m n 1 2.txt,` + TestingCreateContent("m n 10 1 2") + `;
	// `

	input := ""
	destinationInput := ""
	wantedOutcome := ""

	// X move files to existing directories
	// input = input + `
	// 	,2020-11-10T20:40:40Z,txt m e 0.txt,` + contentME110 + `;
	// 	directory 1,2020-11-20T20:40:40Z,txt m e 1.txt,` + contentME111 + `;
	// 	directory 1/directory 2,2020-11-21T20:40:40Z,txt m e 1 2.txt,` + contentME1112 + `;
	// `
	// wantedOutcome = wantedOutcome + `
	// 	2020-11-10,2020-11-10T20:40:40Z,txt m e 0.txt,` + contentME110 + `;
	// 	2020-11-20,2020-11-20T20:40:40Z,txt m e 1.txt,` + contentME111 + `;
	// 	2020-11-21 - 2020-11-22,2020-11-21T20:40:40Z,txt m e 1 2.txt,` + contentME1112 + `;
	// `

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
		2020-11-06 - 2020-11-07,2020-11-06T20:40:40Z,txt 11.txt,` + content112 + `;
		2020-11-06 - 2020-11-07,2020-11-07T20:40:40Z,txt 11 1.txt,` + content1112 + `;
		2020-11-11 - 2020-11-12,2020-11-11T20:40:40Z,txt 11.txt,` + content113 + `;
		2020-11-11 - 2020-11-12,2020-11-12T20:40:40Z,txt 11 1.txt,` + content1113 + `;
		2020-11-16,2020-11-16T20:40:40Z,txt 11.txt,` + content114 + `;
		2020-11-16,2020-11-16T20:40:40Z,txt 11 1.txt,` + content1114 + `;
		2020-11-21 - 2020-11-23,2020-11-21T20:40:40Z,txt 11.txt,` + content115 + `;
		2020-11-21 - 2020-11-23,2020-11-23T20:40:40Z,txt 11 2.txt,` + content116 + `;
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
			destination, _ := utils.TestingCreateFilesAndDirectoriesByOneInput(t, tc.destinationInput)
			defer utils.TestingRemoveDirectoryTree(t, destination)

			wantedOutcomeDestination, _ := utils.TestingCreateFilesAndDirectoriesByOneInput(t, wantedOutcome)
			defer utils.TestingRemoveDirectoryTree(t, wantedOutcomeDestination)

			// act
			if err := filesToDateRangeDirectory(nil, destination); err != nil {
				t.Errorf("filesToDateRangeDirectory error: %v", err)
			}

			if same, err := areDirectoryTreesTheSame(destination, wantedOutcomeDestination); err != nil || !same {
				t.Errorf("compareDirectories error: %v", err)
				t.Errorf("compareDirectories same: %v", same)
			}

			// assert
		})
	}
}
