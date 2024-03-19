package main

import (
	"testing"

	"github.com/markuskreukniet/markus-tools/go/utils/test"
)

func TestFilesToDateRangeDirectory(t *testing.T) {
	// arrange
	testCases := []struct {
		metadata test.TestCaseMetadata
	}{
		{
			metadata: test.TestingCreateTestCaseMetadataWithNameBasicAndWantErrFalse(),
		},
		{
			metadata: test.TestingCreateTestCaseMetadataWithNameEmptyFileSystemNodesAndWantErrFalse(),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.metadata.Name, func(t *testing.T) {
			// arrange and teardown

			// act

			// assert
		})
	}
}
