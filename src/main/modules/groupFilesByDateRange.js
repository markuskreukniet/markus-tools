const fs = require('fs')

export default async function groupFilesByDateRange(filePaths) {
  // path and date created combinations of files
  const pathDateCreatedCombinations = []
  for (const path of filePaths) {
    const stats = fs.statSync(path)
    // Skip zero-byte files and has to be an image file type
    // TODO: duplicate: stats.size > 0
    if (
      stats.size > 0 &&
      (path.endsWith('jpg') ||
        path.endsWith('jpeg') ||
        path.endsWith('png') ||
        path.endsWith('gif'))
    ) {
      pathDateCreatedCombinations.push({ path: path, dateCreated: stats.birthtime })
    }
  }

  // sort combinations
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
