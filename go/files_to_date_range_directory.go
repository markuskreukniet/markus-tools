package main

import (
	"path/filepath"
	"strings"
	"time"

	"github.com/markuskreukniet/markus-tools/go/utils"
)

type fileModificationTime struct {
	path             string
	modificationTime time.Time
}

func filesToDateRangeDirectory(uniqueFileSystemNodes []utils.FileSystemNode, destinationDirectory string) error {
	var paths []string
	if err := utils.AppendFileDetails(
		func(detail utils.FileDetail) {
			paths = append(paths, detail.Path)
		}, uniqueFileSystemNodes, utils.FilesWithoutZeroByteFiles); err != nil {
		return err
	}
	var fileModificationTimes []fileModificationTime
	// TODO: AppendFileDetails should now not look into subdirectories.
	// TODO: utils.Directories is changed from directories
	if err := utils.AppendFileDetails(
		func(detail utils.FileDetail) {
			fileModificationTimes = append(fileModificationTimes, fileModificationTime{
				path:             detail.Path,
				modificationTime: detail.ModificationTime,
			})
		}, []utils.FileSystemNode{{
			Path:        destinationDirectory,
			IsDirectory: true,
		}}, utils.Directories); err != nil {
		return err
	}

	// Remove the directories that are not a date range directory in the destination directory.
	const spacedHyphen = " - "
	for i := 0; i < len(fileModificationTimes); {
		base := filepath.Base(fileModificationTimes[i].path)
		remove := false
		if strings.Contains(base, spacedHyphen) {
			baseParts := strings.Split(base, spacedHyphen)
			if !isValidDateFormat(baseParts[0]) || !isValidDateFormat(baseParts[1]) {
				remove = true
			}
		} else if !isValidDateFormat(base) {
			remove = true
		}
		if remove {
			fileModificationTimes[i] = fileModificationTimes[len(fileModificationTimes)-1]
			fileModificationTimes = fileModificationTimes[:len(fileModificationTimes)-1]
		} else {
			i++
		}
	}

	return nil
}

func isValidDateFormat(date string) bool {
	_, err := time.Parse("2006-01-02", date)
	return err == nil
}
