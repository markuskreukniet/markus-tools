package main

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/markuskreukniet/markus-tools/go/utils"
)

const dateLayout = "2006-01-02" // YYYY-MM-DD

func isValidDateRangeDirectoryName(name string) bool {
	const spacedHyphen = " - "

	parseDate := func(rawDate string) (time.Time, error) {
		date, err := time.Parse(dateLayout, rawDate)
		if err != nil {
			return time.Time{}, err
		}
		return date, nil
	}

	if strings.Contains(name, spacedHyphen) {
		nameParts := strings.Split(name, spacedHyphen)
		firstDate, err := parseDate(nameParts[0])
		if err != nil {
			return false
		}
		secondDate, err := parseDate(nameParts[1])
		if err != nil {
			return false
		}
		daysDifference := secondDate.Sub(firstDate).Hours() / 24
		if daysDifference >= 1 {
			return true
		}
	} else {
		if _, err := parseDate(name); err == nil {
			return true
		}
	}
	return false
}

// The fields in this 'struct' align with the properties of a File instance that we needed in the Kotlin version.
type directoryNameAndFilePath struct {
	directoryName string
	filePath      string
}

func addDirectory(directories *[]string, arg directoryNameAndFilePath) {
	*directories = append(*directories, arg.filePath)
}

func categorizeFilesAndDirectories(
	destinationDirectoryPath string,
) ([]utils.FDateRangeFileInfo, map[string]struct{}, []string, error) {
	var files []utils.FDateRangeFileInfo
	goodDirectoryPaths := make(map[string]struct{})
	var badDirectoryPaths []string

	categorizeInDirectory := func(directoryPaths *[]string, arg directoryNameAndFilePath) {
		if isValidDateRangeDirectoryName(arg.directoryName) {
			goodDirectoryPaths[arg.filePath] = struct{}{}
		} else {
			*directoryPaths = append(*directoryPaths, arg.filePath)
		}
	}

	entries, err := os.ReadDir(destinationDirectoryPath)
	if err != nil {
		return nil, nil, nil, err
	}

	for _, entry := range entries {
		info, err := entry.Info()
		if err != nil {
			return nil, nil, nil, err
		}
		categorize(
			info,
			filepath.Join(destinationDirectoryPath, entry.Name()),
			&files,
			&badDirectoryPaths,
			categorizeInDirectory,
		)
	}

	directories := append([]string{}, badDirectoryPaths...)

	for path := range goodDirectoryPaths {
		directories = append(directories, path)
	}

	for _, path := range directories {
		err := filepath.Walk(path, func(filePath string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}

			if path != filePath {
				categorize(info, filePath, &files, &badDirectoryPaths, addDirectory)
			}

			return nil
		})
		if err != nil {
			return nil, nil, nil, err
		}
	}

	return files, goodDirectoryPaths, badDirectoryPaths, nil
}

func categorize(
	info os.FileInfo,
	filePath string,
	files *[]utils.FDateRangeFileInfo,
	badDirectoryPaths *[]string,
	handler func(*[]string, directoryNameAndFilePath),
) error {
	name := filepath.Base(filePath)

	if info.IsDir() {
		handler(badDirectoryPaths, directoryNameAndFilePath{
			directoryName: name,
			filePath:      filePath,
		})
	} else if info.Mode().IsRegular() {
		size := info.Size()
		if size > 0 {
			*files = append(*files, utils.FDateRangeFileInfo{
				Size:         size,
				Path:         filePath,
				Name:         name,
				TimeModified: info.ModTime(),
			})
		} else {
			// TODO: error
		}
	} else {
		// TODO: error
	}

	return nil
}

func createHandlers(
	destinationDirectoryPath string) []func([]utils.FDateRangeFileInfo, *[]utils.FDateRangeFileInfo,
) []utils.FDateRangeFileInfo {
	appendBadFilesAndReplaceGoodFiles := func(
		badFiles *[]utils.FDateRangeFileInfo, goodFiles *[]utils.FDateRangeFileInfo, file utils.FDateRangeFileInfo,
	) {
		*badFiles = append(*badFiles, *goodFiles...)
		*goodFiles = []utils.FDateRangeFileInfo{file}
	}

	categorizeOnShortestFileNameLength := func(
		files []utils.FDateRangeFileInfo, badFiles *[]utils.FDateRangeFileInfo,
	) []utils.FDateRangeFileInfo {
		getNameLength := func(file utils.FDateRangeFileInfo) int {
			return len(file.Name)
		}

		good := []utils.FDateRangeFileInfo{files[0]}
		var minimumLength = getNameLength(files[0])

		for _, file := range files[1:] {
			nameLength := getNameLength(file)
			if nameLength < minimumLength {
				minimumLength = nameLength
				appendBadFilesAndReplaceGoodFiles(badFiles, &good, file)
			} else if nameLength == minimumLength {
				good = append(good, file)
			} else {
				*badFiles = append(*badFiles, file)
			}
		}

		return good
	}

	categorizeOnValidDateRangeDirectoryName := func(
		files []utils.FDateRangeFileInfo, badFiles *[]utils.FDateRangeFileInfo,
	) []utils.FDateRangeFileInfo {
		var tempGood1Files []utils.FDateRangeFileInfo
		var tempGood2Files []utils.FDateRangeFileInfo
		var tempBadFiles []utils.FDateRangeFileInfo

		for _, file := range files {
			directoryPath := filepath.Dir(file.Path)
			if filepath.Dir(directoryPath) == destinationDirectoryPath {
				if isValidDateRangeDirectoryName(directoryPath) {
					tempGood2Files = append(tempGood2Files, file)
				} else {
					tempGood1Files = append(tempGood1Files, file)
				}
			} else {
				tempBadFiles = append(tempBadFiles, file)
			}
		}

		if len(tempGood2Files) > 0 {
			*badFiles = append(*badFiles, tempGood1Files...)
			*badFiles = append(*badFiles, tempBadFiles...)
			return tempGood2Files
		}

		if len(tempGood1Files) > 0 {
			*badFiles = append(*badFiles, tempBadFiles...)
			return tempGood1Files
		}

		return tempBadFiles
	}

	categorizeOnNewestTimeModified := func(
		files []utils.FDateRangeFileInfo, badFiles *[]utils.FDateRangeFileInfo,
	) []utils.FDateRangeFileInfo {
		good := []utils.FDateRangeFileInfo{files[0]}
		newest := files[0].TimeModified

		for _, file := range files[1:] {
			if file.TimeModified.After(newest) {
				newest = file.TimeModified
				appendBadFilesAndReplaceGoodFiles(badFiles, &good, file)
			} else if file.TimeModified.Equal(newest) {
				good = append(good, file)
			} else {
				*badFiles = append(*badFiles, file)
			}
		}

		return good
	}

	categorizeOnFirstFile := func(
		files []utils.FDateRangeFileInfo, badFiles *[]utils.FDateRangeFileInfo,
	) []utils.FDateRangeFileInfo {
		good := []utils.FDateRangeFileInfo{files[0]}

		*badFiles = append(*badFiles, files[1:]...)

		return good
	}

	return []func([]utils.FDateRangeFileInfo, *[]utils.FDateRangeFileInfo) []utils.FDateRangeFileInfo{
		categorizeOnShortestFileNameLength,
		categorizeOnValidDateRangeDirectoryName,
		categorizeOnNewestTimeModified,
		categorizeOnFirstFile,
	}
}

func deleteDuplicateFiles(files *[]utils.FDateRangeFileInfo, destinationDirectoryPath string) error {
	groups, err := utils.CreateDuplicateFileInfoGroupsByHash(*files, false)
	if err != nil {
		return err
	}

	handlers := createHandlers(destinationDirectoryPath)
	var badFiles []utils.FDateRangeFileInfo

	*files = (*files)[:0]

	for index, group := range groups {
		for _, handler := range handlers {
			// group and groups[index] are different references
			if len(groups[index]) > 1 {
				groups[index] = handler(group, &badFiles)
			} else {
				*files = append(*files, groups[index][0])
				break
			}
		}
	}

	for _, file := range badFiles {
		if err := os.Remove(file.GetPath()); err != nil {
			return err
		}
	}

	return nil
}

func moveFilesAndFilterGoodDirectories(
	files []utils.FDateRangeFileInfo, goodDirectoryPaths *map[string]struct{}, destinationDirectoryPath string,
) error {
	if len(files) == 0 {
		return nil
	}

	sort.Slice(files, func(i, j int) bool {
		return files[i].TimeModified.Before(files[j].TimeModified)
	})

	// A set provides O(1) access time but is unordered, so we cannot reliably retrieve the first value, for example.
	// To maintain order, we also use a slice alongside the set.
	var fileNames map[string]struct{}
	var group []utils.FDateRangeFileInfo

	replaceFileNamesAndGroup := func(file utils.FDateRangeFileInfo) {
		fileNames = map[string]struct{}{file.Name: {}}
		group = []utils.FDateRangeFileInfo{file}
	}

	formatTimeModified := func(file utils.FDateRangeFileInfo) string {
		return file.TimeModified.Format(dateLayout)
	}

	moveFilesToDirectory := func() error {
		firstFile, lastFile := group[0], group[len(group)-1]

		directoryName := formatTimeModified(firstFile)
		if lastFile.TimeModified.Sub(firstFile.TimeModified).Hours() >= 24 {
			directoryName += fmt.Sprintf(" - %s", formatTimeModified(lastFile))
		}

		joinedDirectoryPath := filepath.Join(destinationDirectoryPath, directoryName)

		if _, exists := (*goodDirectoryPaths)[joinedDirectoryPath]; exists {
			delete(*goodDirectoryPaths, joinedDirectoryPath)
		} else {
			if err := utils.CreateDirectory(joinedDirectoryPath); err != nil {
				return err
			}
		}

		for _, file := range group {
			joinedFilePath := filepath.Join(joinedDirectoryPath, file.Name)
			if joinedFilePath != file.Path {
				if err := os.Rename(file.Path, joinedFilePath); err != nil {
					return err
				}
			}
		}

		return nil
	}

	replaceFileNamesAndGroup(files[0])

	for i := 1; i < len(files); i++ {
		lastFile := group[len(group)-1]
		if files[i].TimeModified.Sub(lastFile.TimeModified).Hours() <= 72 {
			if _, exists := fileNames[files[i].Name]; exists {
				extension := filepath.Ext(files[i].Name)
				nameWithoutExtension := strings.TrimSuffix(files[i].Name, extension)
				disambiguationNumber := 2
				for ; disambiguationNumber <= 9; disambiguationNumber++ {
					files[i].Name = fmt.Sprintf("%s %d%s", nameWithoutExtension, disambiguationNumber, extension)
					if _, exists := fileNames[files[i].Name]; !exists {
						break
					}
				}
				if disambiguationNumber == 9 {
					// TODO: error
				}
			}
			fileNames[files[i].Name] = struct{}{}
			group = append(group, files[i])
		} else {
			if err := moveFilesToDirectory(); err != nil {
				return err
			}
			replaceFileNamesAndGroup(files[i])
		}
	}

	if len(group) > 0 {
		if err := moveFilesToDirectory(); err != nil {
			return err
		}
	}

	return nil
}

func filesToDateRangeDirectory(uniqueFileSystemNodes []utils.FileSystemNode, destinationDirectoryPath string) error {
	files, goodDirectoryPaths, badDirectoryPaths, err := categorizeFilesAndDirectories(destinationDirectoryPath)
	if err != nil {
		return err
	}

	for _, node := range uniqueFileSystemNodes {
		info, err := os.Stat(node.Path)
		if err != nil {
			return err
		}
		categorize(info, node.Path, &files, &badDirectoryPaths, addDirectory)
	}

	if err := deleteDuplicateFiles(&files, destinationDirectoryPath); err != nil {
		return err
	}

	if err := moveFilesAndFilterGoodDirectories(files, &goodDirectoryPaths, destinationDirectoryPath); err != nil {
		return err
	}

	// Remove the bad empty directories
	// There is no need to check if the directory exists before attempting removal.
	for i := len(badDirectoryPaths) - 1; i >= 0; i-- {
		if err := os.Remove(badDirectoryPaths[i]); err != nil {
			return err
		}
	}

	for path := range goodDirectoryPaths {
		if err := os.Remove(path); err != nil {
			return err
		}
	}

	return nil
}

func filesToDateRangeDirectoryToJSON(
	uniqueFileSystemNodes []utils.FileSystemNode, destinationDirectoryPath string,
) string {
	if err := filesToDateRangeDirectory(uniqueFileSystemNodes, destinationDirectoryPath); err != nil {
		return errorToJSONFunctionResult(err)
	}
	return defaultJSONFunctionResult()
}
