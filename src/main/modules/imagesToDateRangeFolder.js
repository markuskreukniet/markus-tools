import fs from 'fs'
import isNotAZeroByteFile, { getDirectoryFilePaths } from './filePaths.js'
import { filePathsType } from '../../preload/modules/files'
import { resultStatus, toResultObject } from '../../preload/modules/resultStatus'

// TODO: function looks a lot like duplicateFiles
export default async function imagesToDateRangeFolder(filePaths, path) {
  // TODO: remove the try catch?
  try {
    const directoryFilePathsResult = await getDirectoryFilePaths(
      path,
      false,
      filePathsType.directories
    )

    if (directoryFilePathsResult.status !== resultStatus.ok) {
      return toResultObject(null, directoryFilePathsResult.status, directoryFilePathsResult.message)
    }

    const dateDirectoryFilePaths = getDateSubdirectoryFilePaths(directoryFilePathsResult.result)
    filePaths.push(...dateDirectoryFilePaths)
    const groups = getDateRangeGroups(filePaths)
    groupsToFolders(groups, path)
    // TODO: remove empty folders

    return toResultObject(null, resultStatus.ok)
  } catch (error) {
    return toResultObject(null, resultStatus.errorSystem, error.message)
  }
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
  const filePaths = getFilePaths(path)
  result.push(...filePaths)
  return result
}

function getDateRangeGroups(filePaths) {
  // path and date created combinations of files
  const pathDateCreatedCombinations = []
  for (const path of filePaths) {
    const stats = fs.statSync(path)
    const lowerCasePath = path.toLowerCase()
    if (
      isNotAZeroByteFile(stats) &&
      (lowerCasePath.endsWith('jpg') ||
        lowerCasePath.endsWith('jpeg') ||
        lowerCasePath.endsWith('png') ||
        lowerCasePath.endsWith('gif') ||
        lowerCasePath.endsWith('webp'))
    ) {
      pathDateCreatedCombinations.push({ path, dateCreated: stats.mtime })
    }
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
