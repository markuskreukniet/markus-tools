package main

import (
	"os"
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

const dateFormat = "2006-01-02"

func appendDateRangeDirectoryPathsAndFilePathsTimeModified(dateRangeDirectoryPaths *[]string, filePathsTimeModified *[]filePathTimeModified, filePath string) error {
	*dateRangeDirectoryPaths = append(*dateRangeDirectoryPaths, filePath)
	if err := appendFilePathsTimeModified(filePathsTimeModified, createDirectoryFileSystemNodeInSlice(filePath)); err != nil {
		return err
	}
	return nil
}

func filesToDateRangeDirectory(uniqueFileSystemNodes []utils.FileSystemNode, destinationDirectory string) error {
	var filePathsTimeModified []filePathTimeModified
	if err := appendFilePathsTimeModified(&filePathsTimeModified, uniqueFileSystemNodes); err != nil {
		return err
	}
	const spacedHyphen = " - "
	var dateRangeDirectoryPaths []string
	// TODO: AppendFileDetails should now not look into subdirectories
	// TODO: utils.Directories is changed from directories
	// TODO: AppendFileDetails is probably wrong naming
	// TODO: opErr logic
	var opErr error
	if err := utils.AppendFileDetails(
		func(detail utils.FileDetail) {
			base := filepath.Base(detail.Path)
			if strings.Contains(base, spacedHyphen) {
				baseParts := strings.Split(base, spacedHyphen)
				if isValidDateFormat(baseParts[0]) && isValidDateFormat(baseParts[1]) {
					opErr = appendDateRangeDirectoryPathsAndFilePathsTimeModified(&dateRangeDirectoryPaths, &filePathsTimeModified, detail.Path)
				}
			} else if isValidDateFormat(base) {
				opErr = appendDateRangeDirectoryPathsAndFilePathsTimeModified(&dateRangeDirectoryPaths, &filePathsTimeModified, detail.Path)
			}
		}, createDirectoryFileSystemNodeInSlice(destinationDirectory), utils.Directories); err != nil {
		return err
	}
	if opErr != nil {
		return opErr
	}
	sort.Slice(filePathsTimeModified, func(i, j int) bool {
		return filePathsTimeModified[i].timeModified.Before(filePathsTimeModified[j].timeModified)
	})
	startDateRange := 0
	for i := 0; i < len(filePathsTimeModified)-1; i++ {
		iPlusOne := i + 1
		if filePathsTimeModified[iPlusOne].timeModified.Sub(filePathsTimeModified[i].timeModified).Hours() > 72 {
			var builder strings.Builder
			if _, err := builder.WriteString(toDateFormat(filePathsTimeModified[startDateRange].timeModified)); err != nil {
				return err
			}
			if filePathsTimeModified[startDateRange].timeModified != filePathsTimeModified[i].timeModified {
				if _, err := builder.WriteString(spacedHyphen); err != nil {
					return err
				}
				if _, err := builder.WriteString(toDateFormat(filePathsTimeModified[i].timeModified)); err != nil {
					return err
				}
			}
			subDirectoryPath := filepath.Join(destinationDirectory, builder.String())
			foundIndex := -1
			// TODO: add the missing filePathsTimeModified to the dir and remove subDirectoryPath from dateRangeDirectoryPaths in this loop?
			for j, path := range dateRangeDirectoryPaths {
				if path == subDirectoryPath {
					foundIndex = j
					break
				}
			}
			if foundIndex >= 0 {
				for k := startDateRange; k <= i; k++ {
					fullFilePath := filepath.Join(subDirectoryPath, filepath.Base(filePathsTimeModified[k].filePath))
					if _, err := os.Stat(fullFilePath); err != nil {
						if os.IsNotExist(err) {
							if err := os.Rename(filePathsTimeModified[k].filePath, fullFilePath); err != nil {
								return err
							}
						} else {
							return err
						}
					}
				}
				dateRangeDirectoryPaths[foundIndex] = dateRangeDirectoryPaths[len(dateRangeDirectoryPaths)-1]
				dateRangeDirectoryPaths = dateRangeDirectoryPaths[:len(dateRangeDirectoryPaths)-1]
			} else {
				// if dateRangeDirectoryPaths does not contain subDirectoryPath, make dir with subDirectoryPath
				// add the filePathsTimeModified to the dir from subDirectoryPath
			}
			startDateRange = iPlusOne
		}
	}
	for _, path := range dateRangeDirectoryPaths {
		if err := os.Remove(path); err != nil {
			return err
		}
	}
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
	_, err := time.Parse(dateFormat, date)
	return err == nil
}

func toDateFormat(time time.Time) string {
	return time.Format(dateFormat)
}
