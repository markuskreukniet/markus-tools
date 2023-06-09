import crypto from 'crypto'
import fs from 'fs'
import isNotAZeroByteFile from './filePaths.js'

export default async function duplicateFiles(filePaths) {
  // path and size combinations of files
  const pathSizeCombinations = await Promise.all(
    filePaths.map(async (path) => {
      try {
        const stats = await fs.promises.stat(path)
        // TODO: when getting all files from folder, it should be already isNotAZeroByteFile
        if (isNotAZeroByteFile(stats)) {
          return { path, size: stats.size }
        } else {
          return undefined
        }
      } catch (error) {
        // TODO: error and null
        console.error(error)
        return null
      }
    })
  )

  pathSizeCombinations.sort(compare)

  // duplicates of path and hash combinations
  const duplicates = []
  let lastPushedIndex = -1

  for (let i = 1; i < pathSizeCombinations.length; i++) {
    const combination = pathSizeCombinations[i - 1]
    const combination2 = pathSizeCombinations[i]

    if (combination.size === combination2.size) {
      const hashHex = await getHashHex(combination.path)
      const hashHex2 = await getHashHex(combination2.path)

      if (hashHex === hashHex2) {
        if (lastPushedIndex !== i - 1) {
          duplicates.push({ path: combination.path, hash: hashHex })
        }
        duplicates.push({ path: combination2.path, hash: hashHex2 })
        lastPushedIndex = i
      }
    }
  }

  // result
  if (duplicates.length === 0) {
    return ''
  } else {
    return duplicatesArrayToResultString(duplicates)
  }
}

function getHashHex(path) {
  return new Promise((resolve, reject) => {
    const hash = crypto.createHash('sha1') // SHA1 is faster than MD5
    const stream = fs.createReadStream(path)
    stream.on('error', (err) => reject(err))
    stream.on('data', (chunk) => hash.update(chunk))
    stream.on('end', () => resolve(hash.digest('hex')))
  })
}

// This option is slower when I test it
// 1048576 * 100 // 1 MiB = 1048576 bytes
// const filePartSize = Math.round((1024 * 1024 * 1024) / 10); // Math.round((1 GiB) / 10)
// can be a faster option, but makes hash of a file part
// function getHashHexOfFirstFilePart(path, filePartSize) {
//   return new Promise((resolve, reject) => {
//     const hash = crypto.createHash("sha1"); // SHA1 is faster than MD5
//     const stream = fs.createReadStream(path, {
//       highWaterMark: filePartSize,
//     });
//     stream.on("error", (err) => reject(err));
//     stream.on("data", (chunk) => {
//       resolve(hash.update(chunk).digest("hex"));
//       stream.destroy();
//     });
//   });
// }

function compare(a, b) {
  if (a.size < b.size) {
    return -1
  }
  if (a.size > b.size) {
    return 1
  }
  return 0
}

function duplicatesArrayToResultString(duplicates) {
  let result = duplicates[0].path
  for (let i = 1; i < duplicates.length; i++) {
    let resultPart = `\n${duplicates[i].path}`

    if (duplicates[i].hash !== duplicates[i - 1].hash) {
      resultPart = `\n${resultPart}`
    }

    result += resultPart
  }

  return result
}
