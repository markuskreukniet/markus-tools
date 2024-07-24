package utils

import (
	"path/filepath"
	"sort"
	"strings"
	"testing"
)

func appendFileSystemFilesExtra(filePath string, files *[]FileSystemFileExtra) error {
	handler := func(file FileSystemFile) error {
		hash := ""
		if !file.FileMetadata.IsDirectory {
			var err error
			hash, err = HashFile(file.Path)
			if err != nil {
				return err
			}
		}
		*files = append(*files, CreateFileSystemFileExtra(hash, file))
		return nil
	}

	if err := WalkFilterAndHandleFileSystemFile(filePath, NonZeroByteFilesAndDirectories, AllFiles, handler); err != nil {
		return err
	}

	return nil
}

func areFileSystemFilesExtraIdentical(fileI, fileJ FileSystemFileExtra, filePathI, filePathJ string) (bool, error) {
	// FileMetadata
	// TODO: compare TimeModified
	if fileI.FileSystemFile.FileMetadata.IsDirectory != fileJ.FileSystemFile.FileMetadata.IsDirectory ||
		fileI.FileSystemFile.FileMetadata.Name != fileJ.FileSystemFile.FileMetadata.Name ||
		fileI.FileSystemFile.FileMetadata.Size != fileJ.FileSystemFile.FileMetadata.Size {
		return false, nil
	}

	relativeI, err := filepath.Rel(filePathI, fileI.FileSystemFile.Path)
	if err != nil {
		return false, err
	}

	relativeJ, err := filepath.Rel(filePathJ, fileJ.FileSystemFile.Path)
	if err != nil {
		return false, err
	}

	// FileSystemFile and FileSystemFileExtra
	if relativeI != relativeJ ||
		fileI.FileSystemFile.Data != fileJ.FileSystemFile.Data ||
		fileI.Hash != fileJ.Hash {
		return false, nil
	}

	return true, nil
}

func sortFileSystemFilesExtraOnName(files *[]FileSystemFileExtra) {
	sort.Slice(*files, func(i, j int) bool {
		return (*files)[i].FileSystemFile.FileMetadata.Name < (*files)[j].FileSystemFile.FileMetadata.Name
	})
}

func AreFileTreeDescendantsIdentical(filePathI, filePathJ string) (bool, error) {
	if filePathI == "" || filePathJ == "" {
		return false, nil
	}

	var filesI, filesJ []FileSystemFileExtra

	if err := appendFileSystemFilesExtra(filePathI, &filesI); err != nil {
		return false, err
	}
	if err := appendFileSystemFilesExtra(filePathJ, &filesJ); err != nil {
		return false, err
	}

	length := len(filesI)

	if length != len(filesJ) {
		return false, nil
	}

	filesI[0].FileSystemFile.FileMetadata.Name = ""
	filesJ[0].FileSystemFile.FileMetadata.Name = ""

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
