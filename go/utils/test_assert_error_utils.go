package utils

import (
	"path/filepath"
	"sort"
	"testing"
)

func appendFiles(filePath string, files *[]CompleteFileInfo) error {
	handler := func(file CompleteFileInfo) {
		*files = append(*files, file)
	}

	return WalkFilterAndHandleFileInfoDirectory(filePath, NonZeroByteFilesAndDirectories, AllFiles, handler)
}

func areFilesIdentical(fileI, fileJ CompleteFileInfo, filePathI, filePathJ string) (bool, error) {
	// TODO: compare TimeModified? Or is that part of the hash?
	if fileI.IsDirectory != fileJ.IsDirectory ||
		fileI.Name != fileJ.Name ||
		fileI.Size != fileJ.Size {
		return false, nil
	}

	relativeI, err := filepath.Rel(filePathI, fileI.AbsolutePath)
	if err != nil {
		return false, err
	}

	relativeJ, err := filepath.Rel(filePathJ, fileJ.AbsolutePath)
	if err != nil {
		return false, err
	}

	if relativeI != relativeJ {
		return false, nil
	}

	if !fileI.IsDirectory && !fileJ.IsDirectory {
		hashI, err := CreateFileHash(fileI.AbsolutePath)
		if err != nil {
			return false, err
		}

		hashJ, err := CreateFileHash(fileJ.AbsolutePath)
		if err != nil {
			return false, err
		}

		if hashI != hashJ {
			return false, nil
		}
	}

	return true, nil
}

func AreFileTreeDescendantsIdentical(filePathI, filePathJ string) (bool, error) {
	if filePathI == "" || filePathJ == "" {
		return false, nil
	}

	var filesI, filesJ []CompleteFileInfo

	sortFilesOnName := func(files *[]CompleteFileInfo) {
		sort.Slice(*files, func(i, j int) bool {
			return (*files)[i].Name < (*files)[j].Name
		})
	}

	if err := appendFiles(filePathI, &filesI); err != nil {
		return false, err
	}
	if err := appendFiles(filePathJ, &filesJ); err != nil {
		return false, err
	}

	length := len(filesI)

	if length != len(filesJ) {
		return false, nil
	}

	filesI[0].Name = ""
	filesJ[0].Name = ""

	sortFilesOnName(&filesI)
	sortFilesOnName(&filesJ)

	for i := 1; i < length; i++ {
		areIdentical, err := areFilesIdentical(filesI[i], filesJ[i], filePathI, filePathJ)
		if err != nil {
			return false, err
		}
		if !areIdentical {
			return false, nil
		}
	}

	return true, nil
}

func TMustAssertError(t *testing.T, err error, wantErr bool) {
	t.Helper()

	if (err != nil) != wantErr {
		if wantErr {
			t.Fatalf("want an error, but got nil")
		} else {
			t.Fatalf("did not want an error, but got: %v", err)
		}
	}
}

func TMustAssertEqualStrings(t *testing.T, want string, got string) {
	t.Helper()

	if want != got {
		t.Fatalf("want: %s, got: %s", want, got)
	}
}
