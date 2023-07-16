import fs from 'fs'
import {
  getDirectoryFilePaths,
  getDistinctDirectoryPaths,
  removeEmptyDirectories
} from './filePaths.js'
import { filePathsType, fileType } from '../../preload/modules/files'
import {
  isResultObjectOk,
  isResultObjectPartiallyOk,
  resultStatus,
  toResultObjectWithNullResult,
  toResultObjectWithNullResultAndResultStatusOk,
  toResultObjectWithNullResultByResultObject
} from '../../preload/modules/resultStatus'

// TODO: function looks a lot like duplicateFiles
// TODO: rename resultStatus file
// TODO: check if has access to input en output directory
export default async function imagesToDateRangeFolder(filePaths, outputPath) {
  const inputPath = getSelectedFolderPath(filePaths)

  // TODO: does not work anymore
  // TODO: ResultObject rename to RO
  const imageFilePathsTreeResultObject = await getDirectoryFilePaths(
    inputPath,
    true,
    filePathsType.filesWithoutZeroByteFiles,
    fileType.image
  )
  if (!isResultObjectOk(imageFilePathsTreeResultObject)) {
    return toResultObjectWithNullResultByResultObject(imageFilePathsTreeResultObject)
  }

  const directoryFilePathsResultObject = await getDirectoryFilePaths(
    outputPath,
    false,
    filePathsType.directories
  )
  if (!isResultObjectOk(directoryFilePathsResultObject)) {
    return toResultObjectWithNullResultByResultObject(directoryFilePathsResultObject)
  }

  try {
    const groups = getDateRangeGroups([
      ...imageFilePathsTreeResultObject.result,
      ...getDateSubdirectoryFilePaths(directoryFilePathsResultObject.result)
    ])
    groupsToFolders(groups, outputPath)
  } catch (error) {
    // TODO: use abstraction and also on other places?
    return toResultObjectWithNullResult(resultStatus.errorSystem, error.message)
  }

  const removeEmptyDirectoriesResultObject = await removeEmptyDirectories([
    ...getDistinctDirectoryPaths(imageFilePathsTreeResultObject.result),
    ...directoryFilePathsResultObject.result
  ])
  if (isResultObjectOk(removeEmptyDirectoriesResultObject)) {
    return toResultObjectWithNullResultAndResultStatusOk()
  } else if (isResultObjectPartiallyOk(removeEmptyDirectoriesResultObject)) {
    return toResultObjectWithNullResult(
      resultStatus.partiallyOk,
      removeEmptyDirectoriesResultObject.message
    )
  } else {
    return toResultObjectWithNullResult(
      resultStatus.errorSystem,
      removeEmptyDirectoriesResultObject.message
    )
  }
}

// TODO: remove function
function getSelectedFolderPath(files) {
  const firstFolderPath = getFolderPath(files[0])
  const lastFolderPath = getFolderPath(files[files.length - 1])

  let prefix = ''

  for (let i = 0; i < firstFolderPath.length; i++) {
    if (firstFolderPath[i] === lastFolderPath[i]) {
      prefix += firstFolderPath[i]
    } else {
      break
    }
  }

  return prefix
}

// TODO: remove function
function getFolderPath(file) {
  return file.replace(`\\${file}`, '')
}

function getDateSubdirectoryFilePaths(paths) {
  let result = []
  const separator = ' - '

  for (const path of paths) {
    const pathParts = path.split('\\')
    const lastPathPart = pathParts[pathParts.length - 1]

    if (lastPathPart.includes(separator)) {
      const directoryParts = lastPathPart.split(separator)
      if (isValidDateFormat(directoryParts[0]) && isValidDateFormat(directoryParts[1])) {
        result = addFilePaths(result, path)
      }
    } else if (isValidDateFormat(lastPathPart)) {
      result = addFilePaths(result, path)
    }
  }

  return result
}

function addFilePaths(result, path) {
  return [...result, ...getFilePaths(path)]
}

function getDateRangeGroups(filePaths) {
  // path and date created combinations of files
  const pathDateCreatedCombinations = []
  for (const path of filePaths) {
    const stats = fs.statSync(path)
    pathDateCreatedCombinations.push({ path, dateCreated: stats.mtime })
  }

  pathDateCreatedCombinations.sort(compare)

  const groups = []
  let group = [pathDateCreatedCombinations[0]]
  for (let i = 1; i < pathDateCreatedCombinations.length; i++) {
    const combination = pathDateCreatedCombinations[i - 1]
    const combination2 = pathDateCreatedCombinations[i]

    if (isWithinThreeDays(combination.dateCreated, combination2.dateCreated)) {
      group.push(combination2)
    } else {
      groups.push(group)
      group = [combination2]
    }
  }
  groups.push(group)

  return groups
}

function isWithinThreeDays(date1, date2) {
  const millisecondsDifference = Math.abs(date2.getTime() - date1.getTime())
  const days = millisecondsDifference / (1000 * 3600 * 24)

  return days <= 3
}

function groupsToFolders(groups, path) {
  for (const group of groups) {
    const oldestDate = formatTime(group[0].dateCreated)
    const newestDate = formatTime(group[group.length - 1].dateCreated)

    let subFolderPath = `${path}/${oldestDate}`
    if (oldestDate !== newestDate) {
      subFolderPath = `${subFolderPath} - ${newestDate}`
    }
    if (!fs.existsSync(subFolderPath)) {
      fs.mkdirSync(subFolderPath)
    }
    for (const combination of group) {
      const fileName = combination.path.split('\\').pop().split('/').pop()
      const destinationPath = combinePathAndFile(subFolderPath, fileName)
      fs.rename(combination.path, destinationPath, (err) => {
        if (err) {
          throw err
        }
      })
    }
  }
}

function isValidDateFormat(dateString) {
  const dateParts = dateString.split('-')

  if (dateParts.length !== 3) {
    return false
  }

  const year = parseInt(dateParts[0], 10)
  const month = parseInt(dateParts[1], 10)
  const day = parseInt(dateParts[2], 10)

  return !(
    isNaN(day) ||
    isNaN(month) ||
    isNaN(year) ||
    day < 1 ||
    day > 31 ||
    month < 1 ||
    month > 12 ||
    year < 0
  )
}

function combinePathAndFile(path, file) {
  return `${path}\\${file}`
}

function getFilePaths(path) {
  const filePaths = []

  try {
    const files = fs.readdirSync(path)

    files.forEach((file) => {
      const filePath = combinePathAndFile(path, file)
      const stats = fs.statSync(filePath)

      if (stats.isFile()) {
        filePaths.push(filePath)
      }
    })
  } catch (error) {
    console.error('Error occurred while reading the folder:', error)
  }

  return filePaths
}

function formatTime(time) {
  const date = new Date(time)
  const day = date.getDate()
  // Months are zero-based
  const month = date.getMonth() + 1
  const year = date.getFullYear()

  const formattedDay = addPrefix0IfLessThan10(day)
  const formattedMonth = addPrefix0IfLessThan10(month)

  return `${year}-${formattedMonth}-${formattedDay}`
}

function addPrefix0IfLessThan10(number) {
  return number < 10 ? `0${number}` : number
}

function compare(a, b) {
  if (a.dateCreated < b.dateCreated) {
    return -1
  }
  if (a.dateCreated > b.dateCreated) {
    return 1
  }
  return 0
}
