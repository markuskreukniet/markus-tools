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

const spacedHyphen = " - "
const dateLayout = "2006-01-02" // YYYY-MM-DD

func isValidDateRangeDirectoryName(name string) bool {
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
	} else if isValidDateFormat(name) {
		return true
	}
	return false
}

func isWithin72Hours(olderTime, newerTime time.Time) bool {
	return newerTime.Sub(olderTime).Hours() <= 72
}

func createDirectoryDateRangeName(startTime, endTime time.Time) string {
	start := formatDate(startTime)
	end := formatDate(endTime)

	if start == end {
		return start
	}
	return fmt.Sprintf("%s - %s", start, end)
}

// TODO: naming
type dateRangeArg struct {
	directoryName string
	filePath      string
}

func addDirectory(directories *[]string, arg dateRangeArg) {
	*directories = append(*directories, arg.filePath)
}

func categorizeFilesAndDirectories(destinationDirectory string) ([]utils.DateRangeFileInfo, map[string]struct{}, []string, error) {
	var files []utils.DateRangeFileInfo
	goodDirectoryPaths := make(map[string]struct{})
	var badDirectoryPaths []string

	categorizeInDirectory := func(directoryPaths *[]string, arg dateRangeArg) {
		if isValidDateRangeDirectoryName(arg.directoryName) {
			goodDirectoryPaths[arg.filePath] = struct{}{}
		} else {
			*directoryPaths = append(*directoryPaths, arg.filePath)
		}
	}

	entries, err := os.ReadDir(destinationDirectory)
	if err != nil {
		return nil, nil, nil, err
	}

	for _, entry := range entries {
		info, err := entry.Info()
		if err != nil {
			return nil, nil, nil, err
		}
		categorize(
			info, filepath.Join(destinationDirectory, entry.Name()), &files, &badDirectoryPaths, categorizeInDirectory,
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
	files *[]utils.DateRangeFileInfo,
	badDirectoryPaths *[]string,
	handler func(*[]string, dateRangeArg),
) error {
	name := filepath.Base(filePath)

	if info.IsDir() {
		handler(badDirectoryPaths, dateRangeArg{
			directoryName: name,
			filePath:      filePath,
		})
	} else if info.Mode().IsRegular() {
		size := info.Size()
		if size > 0 {
			*files = append(*files, utils.DateRangeFileInfo{
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

// TODO: Does not work efficient, could be done without making groups?
// garbage collection: length, groups, groupIndex
func moveFilesToDateRangeDirectoriesAndRemoveUsedGoodDirectories(files []utils.FileSystemFile, filePaths []string, filePath string) ([]string, error) {
	length := len(files)

	if length == 0 {
		return filePaths, nil
	}

	sort.Slice(files, func(i, j int) bool {
		return files[i].FileMetadata.TimeModified.Before(files[j].FileMetadata.TimeModified)
	})

	groups := [][]utils.FileSystemFile{{files[0]}}
	groupIndex := 0

	for i := 1; i < length; i++ {
		iMinusOne := i - 1
		if isWithin72Hours(files[iMinusOne].FileMetadata.TimeModified, files[i].FileMetadata.TimeModified) {
			groups[groupIndex] = append(groups[groupIndex], files[i])
		} else {
			groupIndex++
			groups = append(groups, []utils.FileSystemFile{files[i]})
		}
	}

	for _, group := range groups {
		length = len(group)
		lengthMinusOne := length - 1
		var name string
		if group[0].FileMetadata.TimeModified == group[lengthMinusOne].FileMetadata.TimeModified {
			name = formatDate(group[0].FileMetadata.TimeModified)
		} else {
			name = createDirectoryDateRangeName(group[0].FileMetadata.TimeModified, group[lengthMinusOne].FileMetadata.TimeModified)
		}
		directoryFilePath := filepath.Join(filePath, name)
		isDirectoryFound := false
		for j, path := range filePaths {
			if path == directoryFilePath {
				isDirectoryFound = true
				filePaths[j] = filePaths[len(filePaths)-1]
				filePaths = filePaths[:len(filePaths)-1]
				break
			}
		}

		// TODO: should CreateDirectory create a dir with the same rights as parent dir?
		if !isDirectoryFound {
			if err := utils.CreateDirectory(directoryFilePath); err != nil {
				return nil, err
			}
		}

		// TODO clean and make it more efficient
		// add files
		for _, file := range group {
			fullFilePath := filepath.Join(directoryFilePath, file.FileMetadata.Name)
			exists, err := utils.FileExists(fullFilePath)
			if err != nil {
				return nil, err
			}
			if exists {
				// We should always create a hash of the file in the destination folder.
				// Otherwise, we have to loop through all the files to find that file, and that found file might not have a hash yet.
				hash, err := utils.CreateFileHash(fullFilePath)
				if err != nil {
					return nil, err
				}
				if file.FileMetadata.Hash == "" {
					file.FileMetadata.Hash, err = utils.CreateFileHash(file.FileMetadata.Path)
					if err != nil {
						return nil, err
					}
				}
				if hash != file.FileMetadata.Hash {
					extension := filepath.Ext(file.FileMetadata.Name)
					nameWithoutExtension := strings.TrimSuffix(file.FileMetadata.Name, extension)
					fullFilePath = filepath.Join(directoryFilePath, nameWithoutExtension+" 2"+extension)

					if err := os.Rename(file.FileMetadata.Path, fullFilePath); err != nil {
						return nil, err
					}
				}
			} else {
				if err := os.Rename(file.FileMetadata.Path, fullFilePath); err != nil {
					return nil, err
				}
			}
		}
	}

	return filePaths, nil
}

func createHandlers(
	destinationDirectory string) []func([]utils.DateRangeFileInfo, *[]utils.DateRangeFileInfo,
) []utils.DateRangeFileInfo {
	appendBadFilesAndReplaceGoodFiles := func(
		badFiles *[]utils.DateRangeFileInfo, goodFiles *[]utils.DateRangeFileInfo, file utils.DateRangeFileInfo,
	) {
		*badFiles = append(*badFiles, *goodFiles...)
		*goodFiles = []utils.DateRangeFileInfo{file}
	}

	categorizeOnShortestFileNameLength := func(
		files []utils.DateRangeFileInfo, badFiles *[]utils.DateRangeFileInfo,
	) []utils.DateRangeFileInfo {
		getNameLength := func(file utils.DateRangeFileInfo) int {
			return len(file.Name)
		}

		good := []utils.DateRangeFileInfo{files[0]}
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
		files []utils.DateRangeFileInfo, badFiles *[]utils.DateRangeFileInfo,
	) []utils.DateRangeFileInfo {
		var tempGood1Files []utils.DateRangeFileInfo
		var tempGood2Files []utils.DateRangeFileInfo
		var tempBadFiles []utils.DateRangeFileInfo

		for _, file := range files {
			directoryPath := filepath.Dir(file.Path)
			if filepath.Dir(directoryPath) == destinationDirectory {
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
		files []utils.DateRangeFileInfo, badFiles *[]utils.DateRangeFileInfo,
	) []utils.DateRangeFileInfo {
		good := []utils.DateRangeFileInfo{files[0]}
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
		files []utils.DateRangeFileInfo, badFiles *[]utils.DateRangeFileInfo,
	) []utils.DateRangeFileInfo {
		good := []utils.DateRangeFileInfo{files[0]}

		*badFiles = append(*badFiles, files[1:]...)

		return good
	}

	return []func([]utils.DateRangeFileInfo, *[]utils.DateRangeFileInfo) []utils.DateRangeFileInfo{
		categorizeOnShortestFileNameLength,
		categorizeOnValidDateRangeDirectoryName,
		categorizeOnNewestTimeModified,
		categorizeOnFirstFile,
	}
}

func deleteDuplicateFiles(files *[]utils.DateRangeFileInfo, destinationDirectory string) error {
	groups, err := utils.CreateDuplicateFileInfoGroupsByHash(*files, false)
	if err != nil {
		return err
	}

	handlers := createHandlers(destinationDirectory)
	var badFiles []utils.DateRangeFileInfo

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

// TODO: rename destinationDirectory to destinationDirectoryPath on other places
func moveFilesAndFilterGoodDirectories(
	files []utils.DateRangeFileInfo, goodDirectoryPaths *map[string]struct{}, destinationDirectoryPath string,
) error {
	if len(files) == 0 {
		return nil
	}

	sort.Slice(files, func(i, j int) bool {
		return files[i].TimeModified.Before(files[j].TimeModified)
	})

	// TODO: also do this set logic in Kotlin
	// TODO: Grammar: a set is not ordered, so we can´t take the first value for example. But it has O(1) access time. So we need a set and a slice.
	var fileNames map[string]struct{}
	var group []utils.DateRangeFileInfo

	replaceFileNamesAndGroup := func(file utils.DateRangeFileInfo) {
		fileNames = map[string]struct{}{file.Name: {}}
		group = []utils.DateRangeFileInfo{file}
	}

	formatTimeModified := func(file utils.DateRangeFileInfo) string {
		return file.TimeModified.Format(dateLayout)
	}

	moveFilesToDirectory := func() error {
		firstFile := group[0]
		lastFile := group[len(group)-1]

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

	// TODO: search for i := 1 for range files[1:]
	// TODO: duplicate code does not work on Linux
	// TODO: also fix in Kotlin code

	for i := 1; i < len(files); i++ {
		lastFile := group[len(group)-1]
		if files[i].TimeModified.Sub(lastFile.TimeModified).Hours() <= 72 {
			if _, exists := fileNames[files[i].Name]; exists {
				extension := filepath.Ext(files[i].Name)
				name := strings.TrimSuffix(files[i].Name, extension) + " 2" + extension // TODO: with " 2" might also exists
				files[i].Name = name
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

func filesToDateRangeDirectory(uniqueFileSystemNodes []utils.FileSystemNode, destinationDirectory string) error {
	files, goodDirectoryPaths, badDirectoryPaths, err := categorizeFilesAndDirectories(destinationDirectory)
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

	if err := deleteDuplicateFiles(&files, destinationDirectory); err != nil {
		return err
	}

	if err := moveFilesAndFilterGoodDirectories(files, &goodDirectoryPaths, destinationDirectory); err != nil {
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

func formatDate(time time.Time) string {
	return time.Format(dateLayout)
}

func parseDate(rawDate string) (time.Time, error) {
	date, err := time.Parse(dateLayout, rawDate)
	if err != nil {
		return time.Time{}, err
	}
	return date, nil
}

func isValidDateFormat(rawDate string) bool {
	_, err := parseDate(rawDate)
	return err == nil
}
