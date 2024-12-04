package utils

import (
	"path/filepath"
	"sort"
	"strings"
	"testing"
)

func appendFiles(filePath string, files *[]CompleteFileInfo) error {
	handler := func(file CompleteFileInfo) {
		*files = append(*files, file)
	}

	return WalkFilterAndHandleFileInfoDirectory(filePath, NonZeroByteFilesAndDirectories, AllFiles, handler)
}

func areFilesIdentical(fileI, fileJ CompleteFileInfo, filePathI, filePathJ string) (bool, error) {
	// TODO: compare TimeModified
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

func sortFilesOnName(files *[]CompleteFileInfo) {
	sort.Slice(*files, func(i, j int) bool {
		return (*files)[i].Name < (*files)[j].Name
	})
}

func AreFileTreeDescendantsIdentical(filePathI, filePathJ string) (bool, error) {
	if filePathI == "" || filePathJ == "" {
		return false, nil
	}

	var filesI, filesJ []CompleteFileInfo

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

func TestingAssertErrorToWantErrorAndOutcomeToBuilderString(t *testing.T, err error, wantErr bool, builder strings.Builder, got string) {
	t.Helper()
	TestingAssertErrorToWantError(t, err, wantErr)
	TestingAssertEqualStrings(t, builder.String(), got)
}

func TestingAssertErrorToWantError(t *testing.T, err error, wantErr bool) {
	t.Helper()
	if (err != nil) != wantErr {
		t.Errorf("want error: %v, got error: %v", wantErr, err)
	}
}

func TestingAssertEqualStrings(t *testing.T, want string, got string) {
	t.Helper()
	if want != got {
		t.Errorf("want: %s, got: %s", want, got)
	}
}
