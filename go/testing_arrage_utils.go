package main

type TestCaseMetadata struct {
	Name    string
	WantErr bool
}

func testingCreateTestCaseMetadata(name string, wantErr bool) TestCaseMetadata {
	return TestCaseMetadata{
		Name:    name,
		WantErr: wantErr,
	}
}

func testingCreateTestCaseMetadataWithWantErrTrue(name string) TestCaseMetadata {
	return testingCreateTestCaseMetadata(name, true)
}

func testingCreateTestCaseMetadataWithNameBasicAndWantErrFalse() TestCaseMetadata {
	return testingCreateTestCaseMetadata("Basic", false)
}

func testingCreateTestCaseMetadataWithNameEmptyFileSystemNodesAndWantErrFalse() TestCaseMetadata {
	return testingCreateTestCaseMetadata("Empty FileSystemNodes", false)
}
