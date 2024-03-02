package main

import (
	"path/filepath"
	"strings"
	"time"

	"github.com/markuskreukniet/markus-tools/go/utils"
)

type fileTimeModified struct {
	path         string
	timeModified time.Time
}

func filesToDateRangeDirectory(uniqueFileSystemNodes []utils.FileSystemNode, destinationDirectory string) error {
	var filesTimeModified []fileTimeModified
	if err := utils.AppendFileDetails(
		func(detail utils.FileDetail) {
			filesTimeModified = append(filesTimeModified, fileTimeModified{
				path:         detail.Path,
				timeModified: detail.ModificationTime,
			})
		}, uniqueFileSystemNodes, utils.FilesWithoutZeroByteFiles); err != nil {
		return err
	}
	var dateRangeDirectoryPaths []string
	// TODO: AppendFileDetails should now not look into subdirectories???
	// TODO: utils.Directories is changed from directories
	if err := utils.AppendFileDetails(
		func(detail utils.FileDetail) {
			dateRangeDirectoryPaths = append(dateRangeDirectoryPaths, detail.Path)
		}, []utils.FileSystemNode{{
			Path:        destinationDirectory,
			IsDirectory: true,
		}}, utils.Directories); err != nil {
		return err
	}

	// Remove the directories from the slice that are not a date range directory in the destination directory.
	const spacedHyphen = " - "
	for i := 0; i < len(dateRangeDirectoryPaths); {
		base := filepath.Base(dateRangeDirectoryPaths[i])
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
			dateRangeDirectoryPaths[i] = dateRangeDirectoryPaths[len(dateRangeDirectoryPaths)-1]
			dateRangeDirectoryPaths = dateRangeDirectoryPaths[:len(dateRangeDirectoryPaths)-1]
		} else {
			// TODO: get all files of that dir
			i++
		}
	}

	// TODO: get all files of the dateRangeDirectories and add them to filesTimeModified.
	// Some dirs might be empty, which we should delete, or add to empty dir slice
	// for _, directory := range dateRangeDirectories {

	// }

	// TODO: order filesTimeModified
	// TODO: filesTimeModified to groups
	// TODO: remove dateRangeDirectoryPaths that have a corresponding group
	// TODO: groups to directories
	// TODO: groups to files in directories
	// TODO: remove all directories from dateRangeDirectoryPaths

	// might have already overlapping date range dirs

	return nil
}

func isValidDateFormat(date string) bool {
	_, err := time.Parse("2006-01-02", date)
	return err == nil
}
