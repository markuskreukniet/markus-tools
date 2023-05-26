const fs = require('fs')
import isNotAZeroByteFile from './fileHelper.js'

export default async function groupFilesByDateRange(filePaths) {
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
      pathDateCreatedCombinations.push({ path: path, dateCreated: stats.birthtime })
    }
  }

  // TODO: sort duplicate
  pathDateCreatedCombinations.sort(compare)
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
