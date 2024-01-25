package main

import "path/filepath"

type TestCaseMetadata struct {
	Name    string
	WantErr bool
}

type FileSystemPathEndParts struct {
	DirectoryPathEndParts []string
	FilePathEndParts      []string
}

// TODO: there are no files in the root (temp dir)
var (
	directoryEmpty               = "directory empty"
	directory1                   = "directory 1"
	directory2                   = "directory 2"
	directory2WithDirectoryEmpty = filepath.Join(directory2, directoryEmpty)
	directory2WithDirectory3     = filepath.Join(directory2, "directory 3")
	directory2WithDirectory4     = filepath.Join(directory2, "directory 4")

	txtFile1 = filepath.Join(directory1, "file 1.txt")
	txtFile2 = filepath.Join(directory1, "file 2.txt")
	txtFile3 = filepath.Join(directory2WithDirectory3, "file 3.txt")
	txtFile4 = filepath.Join(directory2WithDirectory3, "file 4.txt")
	txtFile5 = filepath.Join(directory2WithDirectory3, "file 5.txt")
	txtFile6 = filepath.Join(directory2WithDirectory4, "file 6.txt")

	jpgFile4 = filepath.Join(directory1, "file 4.jpg")

	txtFileNonExistent1 = "non existent 1.txt"

	emptyPathEndParts []string

	emptyFileSystemPathEndParts = FileSystemPathEndParts{
		DirectoryPathEndParts: emptyPathEndParts,
		FilePathEndParts:      emptyPathEndParts,
	}
)

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
