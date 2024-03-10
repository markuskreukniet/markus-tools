package main

import (
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/markuskreukniet/markus-tools/go/utils"
)

type filePathTimeModified struct {
	filePath     string
	timeModified time.Time
}

func filesToDateRangeDirectory(uniqueFileSystemNodes []utils.FileSystemNode, destinationDirectory string) error {
	var filePathsTimeModified []filePathTimeModified
	if err := appendFilePathsTimeModified(&filePathsTimeModified, uniqueFileSystemNodes); err != nil {
		return err
	}
	var dateRangeDirectoryPaths []string
	// TODO: AppendFileDetails should now not look into subdirectories
	// TODO: utils.Directories is changed from directories
	if err := utils.AppendFileDetails(
		func(detail utils.FileDetail) {
			dateRangeDirectoryPaths = append(dateRangeDirectoryPaths, detail.Path)
		}, createDirectoryFileSystemNodeInSlice(destinationDirectory), utils.Directories); err != nil {
		return err
	}

	// Remove the directories from the slice that are not a date range directory in the destination directory.
	// And append the filePathsTimeModified from a date range directory.
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
			if err := appendFilePathsTimeModified(&filePathsTimeModified, createDirectoryFileSystemNodeInSlice(dateRangeDirectoryPaths[i])); err != nil {
				return err
			}
			i++
		}
	}

	// TODO: if length 0 stop?

	sort.Slice(filePathsTimeModified, func(i, j int) bool {
		return filePathsTimeModified[i].timeModified.Before(filePathsTimeModified[j].timeModified)
	})
	var dateRanges [][]filePathTimeModified
	dateRange := []filePathTimeModified{filePathsTimeModified[0]}
	for i := 0; i < len(filePathsTimeModified); i++ {
		if filePathsTimeModified[i].timeModified.Sub(filePathsTimeModified[i-1].timeModified).Hours() <= 72 {
			dateRange = append(dateRange, filePathsTimeModified[i])
		} else {
			dateRanges = append(dateRanges, dateRange)
			dateRange = []filePathTimeModified{filePathsTimeModified[i]}
		}
	}

	// TODO: filePathsTimeModified to groups
	// TODO: remove dateRangeDirectoryPaths that have a corresponding group
	// TODO: groups to directories
	// TODO: groups to files in directories
	// TODO: remove all directories from dateRangeDirectoryPaths

	// might have already overlapping date range dirs

	return nil
}

func createDirectoryFileSystemNodeInSlice(path string) []utils.FileSystemNode {
	return []utils.FileSystemNode{{
		Path:        path,
		IsDirectory: true,
	}}
}

func appendFilePathsTimeModified(filePathsTimeModified *[]filePathTimeModified, uniqueFileSystemNodes []utils.FileSystemNode) error {
	return utils.AppendFileDetails(
		func(detail utils.FileDetail) {
			*filePathsTimeModified = append(*filePathsTimeModified, filePathTimeModified{
				filePath:     detail.Path,
				timeModified: detail.ModificationTime,
			})
		}, uniqueFileSystemNodes, utils.FilesWithoutZeroByteFiles)
}

func isValidDateFormat(date string) bool {
	_, err := time.Parse("2006-01-02", date)
	return err == nil
}
