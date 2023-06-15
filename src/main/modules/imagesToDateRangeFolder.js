import fs from 'fs'
import isNotAZeroByteFile from './fileHelper.js'

// TODO: function looks a lot like duplicateFiles
export default async function imagesToDateRangeFolder(filePaths) {
  // TODO: also in duplicateFiles and check return type
  if (filePaths.length < 2) {
    return false
  }

  const groups = getDateRangeGroups(filePaths)
  groupsToFolders(groups, '')

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

function groupsToFolders(groups, path) {}

function compare(a, b) {
  if (a.dateCreated < b.dateCreated) {
    return -1
  }
  if (a.dateCreated > b.dateCreated) {
    return 1
  }
  return 0
}
