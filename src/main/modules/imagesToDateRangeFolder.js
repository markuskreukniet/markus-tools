import fs from 'fs'
import isNotAZeroByteFile from './fileHelper.js'

// TODO: function looks a lot like duplicateFiles
export default async function imagesToDateRangeFolder(filePaths, path) {
  // TODO: also in duplicateFiles and check return type
  if (filePaths.length < 2) {
    return false
  }

  const groups = getDateRangeGroups(filePaths)
  groupsToFolders(groups, path)

  return true
}

function getDateRangeGroups(filePaths) {
  // path and date created combinations of files
  const pathDateCreatedCombinations = []
  for (const path of filePaths) {
    const stats = fs.statSync(path)
    if (
      isNotAZeroByteFile(stats) &&
      (path.endsWith('jpg') ||
        path.endsWith('jpeg') ||
        path.endsWith('png') ||
        path.endsWith('gif'))
    ) {
      pathDateCreatedCombinations.push({ path, dateCreated: stats.birthtime })
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

function getDirectoryDates(directories) {
  const result = []
  const separator = ' - '

  for (const directory of directories) {
    if (directory.includes(separator)) {
      const directoryParts = directory.split(separator)
      if (isValidDateFormat(directoryParts[0]) && isValidDateFormat(directoryParts[1])) {
        const date1 = formattedDateToDate(directoryParts[0])
        const date2 = formattedDateToDate(directoryParts[1])
        result.push([date1, date2])
      }
    } else if (isValidDateFormat(directory)) {
      const date = formattedDateToDate(directory)
      result.push([date])
    }
  }

  return result
}

function setOldestDate(oldestDate, dates) {
  let result = oldestDate

  for (let i = dates.length - 1; i >= 0; i--) {
    if (dates[i] < result && isWithinThreeDays(dates[i], result)) {
      result = dates[i]
    }
  }

  return result
}

function setNewestDate(newestDate, dates) {
  let result = newestDate

  for (const date of dates) {
    if (date > result && isWithinThreeDays(date, result)) {
      result = date
    }
  }

  return result
}

function getSubFolderPath(path, oldestDate, newestDate) {
  let subFolderPathFiles = `${path}/${oldestDate}`
  if (oldestDate !== newestDate) {
    subFolderPathFiles = `${subFolderPathFiles} - ${newestDate}`
  }
  return subFolderPathFiles
}

function groupsToFolders(groups, path) {
  const directories = getSubdirectories(path)
  const directoryDates = getDirectoryDates(directories)

  for (const group of groups) {
    let oldestDate = formatBirthtime(group[0].dateCreated)
    let newestDate = formatBirthtime(group[group.length - 1].dateCreated)

    let subFolderPathFiles = getSubFolderPath(path, oldestDate, newestDate)

    for (const dates of directoryDates) {
      oldestDate = setOldestDate(oldestDate, dates)
      newestDate = setNewestDate(oldestDate, dates)
    }

    let subFolderPathFilesAndFolder = getSubFolderPath(path, oldestDate, newestDate)

    if (!fs.existsSync(subFolderPathFiles)) {
      fs.mkdirSync(subFolderPathFiles)
    }

    for (const combination of group) {
      const fileName = combination.path.split('\\').pop().split('/').pop()

      fs.copyFile(combination.path, `${subFolderPathFiles}/${fileName}`, (err) => {
        if (err) {
          throw err
        }
      })
    }
  }
}

function isBetweenDate(firstDate, secondDate, thirdDate) {
  const newFirstDate = formattedDateToDate(firstDate)
  const newSecondDate = formattedDateToDate(secondDate)
  const newThirdDate = formattedDateToDate(thirdDate)

  return newThirdDate > newFirstDate && newThirdDate < newSecondDate
}

function formattedDateToDate(formattedDate) {
  const formattedDateParts = formattedDate.split('-')

  const day = parseInt(formattedDateParts[0], 10)
  // Months are zero-based, therefore - 1
  const month = parseInt(formattedDateParts[1], 10) - 1
  const year = parseInt(formattedDateParts[2], 10)

  return new Date(year, month, day)
}

function isValidDateFormat(dateString) {
  const dateParts = dateString.split('-')

  if (dateParts.length !== 3) {
    return false
  }

  const day = parseInt(dateParts[0], 10)
  const month = parseInt(dateParts[1], 10)
  const year = parseInt(dateParts[2], 10)

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

function getSubdirectories(path) {
  const subdirectories = []

  try {
    const files = fs.readdirSync(path)

    files.forEach((file) => {
      const stats = fs.statSync(`${path}/${file}`)

      if (stats.isDirectory()) {
        subdirectories.push(file)
      }
    })
  } catch (error) {
    console.error('Error occurred while reading the folder:', error)
  }

  return subdirectories
}

function formatBirthtime(birthtime) {
  const date = new Date(birthtime)
  const day = date.getDate()
  // Months are zero-based
  const month = date.getMonth() + 1
  const year = date.getFullYear()

  const formattedDay = addPrefix0IfLessThan10(day)
  const formattedMonth = addPrefix0IfLessThan10(month)

  return `${formattedDay}-${formattedMonth}-${year}`
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
