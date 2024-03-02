package main

import (
	"path/filepath"
	"strings"
	"time"

	"github.com/markuskreukniet/markus-tools/go/utils"
)

type fileModifiedFormattedDate struct {
	path                  string
	modifiedFormattedDate string
}

type dateRangeDirectory struct {
	path      string
	startDate string
	endDate   string
}

const dateFormat = "2006-01-02"

func filesToDateRangeDirectory(uniqueFileSystemNodes []utils.FileSystemNode, destinationDirectory string) error {
	var fileModifiedFormattedDates []fileModifiedFormattedDate
	if err := utils.AppendFileDetails(
		func(detail utils.FileDetail) {
			fileModifiedFormattedDates = append(fileModifiedFormattedDates, fileModifiedFormattedDate{
				path:                  detail.Path,
				modifiedFormattedDate: detail.ModificationTime.Format(dateFormat),
			})
		}, uniqueFileSystemNodes, utils.FilesWithoutZeroByteFiles); err != nil {
		return err
	}
	var dateRangeDirectories []dateRangeDirectory
	// TODO: AppendFileDetails should now not look into subdirectories.
	// TODO: utils.Directories is changed from directories
	if err := utils.AppendFileDetails(
		func(detail utils.FileDetail) {
			dateRangeDirectories = append(dateRangeDirectories, dateRangeDirectory{
				path:      detail.Path,
				startDate: "",
				endDate:   "",
			})
		}, []utils.FileSystemNode{{
			Path:        destinationDirectory,
			IsDirectory: true,
		}}, utils.Directories); err != nil {
		return err
	}

	// Remove the directories from the slice that are not a date range directory in the destination directory.
	const spacedHyphen = " - "
	for i := 0; i < len(dateRangeDirectories); {
		base := filepath.Base(dateRangeDirectories[i].path)
		startDate := ""
		endDate := ""
		remove := false
		if strings.Contains(base, spacedHyphen) {
			baseParts := strings.Split(base, spacedHyphen)
			if isValidDateFormat(baseParts[0]) && isValidDateFormat(baseParts[1]) {
				startDate = baseParts[0]
				endDate = baseParts[1]
			} else {
				remove = true
			}
		} else if isValidDateFormat(base) {
			startDate = base
		} else {
			remove = true
		}
		if remove {
			dateRangeDirectories[i] = dateRangeDirectories[len(dateRangeDirectories)-1]
			dateRangeDirectories = dateRangeDirectories[:len(dateRangeDirectories)-1]
		} else {
			dateRangeDirectories[i].startDate = startDate
			dateRangeDirectories[i].endDate = endDate
			i++
		}
	}

	//
	// for _, time := range fileModificationTimes {
	// 	for _, directory := range dateRangeDirectories {
	// 		if directory.endDate == "" && time. {

	// 		}
	// 	}
	// }

	return nil
}

func isValidDateFormat(date string) bool {
	_, err := time.Parse(dateFormat, date)
	return err == nil
}
