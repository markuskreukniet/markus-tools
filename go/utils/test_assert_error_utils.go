package utils

import (
	"path/filepath"
	"sort"
	"strings"
	"testing"
)

func appendFileSystemFiles(filePath string, files *[]FileSystemFile) error {
	handler := func(file FileSystemFile) error {
		if !file.FileMetadata.IsDirectory {
			var err error
			file.FileMetadata.Hash, err = CreateFileHash(file.FileMetadata.Path)
			if err != nil {
				return err
			}
		}
		*files = append(*files, file)
		return nil
	}

	if err := WalkFilterAndHandleFileSystemFile(filePath, NonZeroByteFilesAndDirectories, AllFiles, handler); err != nil {
		return err
	}

	return nil
}

func areFileSystemFilesExtraIdentical(fileI, fileJ FileSystemFile, filePathI, filePathJ string) (bool, error) {
	// FileMetadata
	// TODO: compare TimeModified
	if fileI.FileMetadata.IsDirectory != fileJ.FileMetadata.IsDirectory ||
		fileI.FileMetadata.Name != fileJ.FileMetadata.Name ||
		fileI.FileMetadata.Size != fileJ.FileMetadata.Size ||
		fileI.FileMetadata.Hash != fileJ.FileMetadata.Hash {
		return false, nil
	}

	relativeI, err := filepath.Rel(filePathI, fileI.FileMetadata.Path)
	if err != nil {
		return false, err
	}

	relativeJ, err := filepath.Rel(filePathJ, fileJ.FileMetadata.Path)
	if err != nil {
		return false, err
	}

	// FileSystemFile
	if relativeI != relativeJ ||
		fileI.Data != fileJ.Data {
		return false, nil
	}

	return true, nil
}

func sortFileSystemFilesExtraOnName(files *[]FileSystemFile) {
	sort.Slice(*files, func(i, j int) bool {
		return (*files)[i].FileMetadata.Name < (*files)[j].FileMetadata.Name
	})
}

func AreFileTreeDescendantsIdentical(filePathI, filePathJ string) (bool, error) {
	if filePathI == "" || filePathJ == "" {
		return false, nil
	}

	var filesI, filesJ []FileSystemFile

	if err := appendFileSystemFiles(filePathI, &filesI); err != nil {
		return false, err
	}
	if err := appendFileSystemFiles(filePathJ, &filesJ); err != nil {
		return false, err
	}

	length := len(filesI)

	if length != len(filesJ) {
		return false, nil
	}

	filesI[0].FileMetadata.Name = ""
	filesJ[0].FileMetadata.Name = ""

	sortFileSystemFilesExtraOnName(&filesI)
	sortFileSystemFilesExtraOnName(&filesJ)

	for i := 1; i < length; i++ {
		areIdentical, err := areFileSystemFilesExtraIdentical(filesI[i], filesJ[i], filePathI, filePathJ)
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
