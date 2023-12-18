package main

import "path/filepath"

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
)
