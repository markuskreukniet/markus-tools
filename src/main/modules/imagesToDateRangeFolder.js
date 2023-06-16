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

function groupsToFolders(groups, path) {
  // const directories = getSubdirectories(path)

  for (const group of groups) {
    const oldestDate = formatBirthtime(group[0].dateCreated)
    const newestDate = formatBirthtime(group[group.length - 1].dateCreated)

    // for (const directory of directories) {
    // }

    let subFolderPath = `${path}/${oldestDate}`
    if (oldestDate !== newestDate) {
      subFolderPath = `${subFolderPath} - ${newestDate}`
    }

    if (!fs.existsSync(subFolderPath)) {
      fs.mkdirSync(subFolderPath)
    }

    for (const combination of group) {
      const fileName = combination.path.split('\\').pop().split('/').pop()

      fs.copyFile(combination.path, `${subFolderPath}/${fileName}`, (err) => {
        if (err) {
          throw err
        }
      })
    }
  }
}

function isBetweenDate(firstDate, secondDate, thirdDate) {
  const firstParts = firstDate.split('-')
  const secondParts = secondDate.split('-')
  const thirdParts = thirdDate.split('-')

  const firstDay = parseInt(firstParts[0], 10)
  // Months are zero-based
  const firstMonth = parseInt(firstParts[1], 10) - 1
  const firstYear = parseInt(firstParts[2], 10)

  const secondDay = parseInt(secondParts[0], 10)
  // Months are zero-based
  const secondMonth = parseInt(secondParts[1], 10) - 1
  const secondYear = parseInt(secondParts[2], 10)

  const thirdDay = parseInt(thirdParts[0], 10)
  // Months are zero-based
  const thirdMonth = parseInt(thirdParts[1], 10) - 1
  const thirdYear = parseInt(thirdParts[2], 10)

  const newFirstDate = new Date(firstYear, firstMonth, firstDay)
  const newSecondDate = new Date(secondYear, secondMonth, secondDay)
  const newThirdDate = new Date(thirdYear, thirdMonth, thirdDay)

  return newThirdDate > newFirstDate && newThirdDate < newSecondDate
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
