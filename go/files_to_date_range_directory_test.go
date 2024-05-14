package main

import (
	"log"
	"testing"

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

func TestFilesToDateRangeDirectory(t *testing.T) {
	// arrange
	time0240 := "2020-02-20T20:40:40Z"
	time0241 := "2020-02-20T20:40:41Z"
	time05 := "2020-05-20T20:40:40Z"
	input := `
		,,txt 0.txt,;
		empty,,,;
		directory 1,` + time0241 + `,txt 1.txt,;
		directory 1,,jpg 1.jpg,;
		directory 2/directory 3,` + time05 + `,txt 2 3.txt,;
	`
	destinationInput := `
		,,txt 0.txt,;
		2020-01-20,,,;
		2020-02-20,` + time0240 + `,txt 02.txt,;
		2020-03-20,,,;
		2020-04-20 - 2020-04-21,,,;
		2020-05-20 - 2020-05-21,` + time05 + `,txt 05.txt,;
		2020-06-20 - 2020-06-21,,,;
	`
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
			directories, _ := utils.TestingCreateFilesAndDirectoriesByMultipleInputs(t, tc.testCaseInput.Input)
			defer utils.TestingRemoveDirectoryTrees(t, directories)
			directory, _ := utils.TestingCreateFilesAndDirectoriesByOneInput(t, tc.destinationInput)
			defer utils.TestingRemoveDirectoryTree(t, directory)

			log.Println(directories)
			log.Println(directory)
		})
	}

	// act
	// assert
}
