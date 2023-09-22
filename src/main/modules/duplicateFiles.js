import crypto from 'crypto'
import { filePathObjectsToFileObjects, getReadFileHandle } from './filePaths.js'
import {
  isResultObjectOk,
  toResultObjectWithNullResultByResultObject
} from '../../preload/modules/resultStatus'

// TODO: return RO
// TODO: has many similarities with imagesToDateRangeFolder.js
export default async function duplicateFiles(filePathObjects) {
  // TODO: filePathObjects > 0
  // TODO: tree optional
  const filePathObjectsToFileObjectsRO = await filePathObjectsToFileObjects(filePathObjects, true)
  if (!isResultObjectOk(filePathObjectsToFileObjectsRO)) {
    return toResultObjectWithNullResultByResultObject(filePathObjectsToFileObjectsRO)
  }

  const fileObjects = filePathObjectsToFileObjectsRO.result
  fileObjects.sort(compare)

  const duplicates = []
  let lastPushedIndex = -1

  for (let i = 1; i < fileObjects.length; i++) {
    const fileObject = fileObjects[i - 1]
    const fileObject2 = fileObjects[i]

    if (fileObject.size === fileObject2.size) {
      const fileHash = await getFileHash(fileObject.path)
      const fileHash2 = await getFileHash(fileObject2.path)

      if (fileHash === fileHash2) {
        if (lastPushedIndex !== i - 1) {
          // TODO: create object with function
          duplicates.push({ path: fileObject.path, hash: fileHash })
        }
        duplicates.push({ path: fileObject2.path, hash: fileHash2 })
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

async function getFileHash(filePath) {
  const fileHandleRO = await getReadFileHandle(filePath)
  if (!isResultObjectOk(fileHandleRO)) {
    return toResultObjectWithNullResultByResultObject(fileHandleRO)
  }

  let readStream = null
  try {
    readStream = fileHandleRO.result.createReadStream()
  } catch {
    // TODO:
  }

  // TODO: change to try catch?
  let fileHash = ''
  await getFileHashByReadStream(readStream)
    .then((hash) => {
      fileHash = hash
    })
    .catch((error) => {
      console.error('Error:', error) // TODO:
    })

  fileHandleRO.result.close()
  return fileHash
}

function getFileHashByReadStream(readStream) {
  // sha256 is generally faster and more secure than SHA1
  // SHA1 is generally faster and more secure than MD5
  const hash = crypto.createHash('sha256')

  return new Promise((resolve, reject) => {
    readStream.on('data', (chunk) => hash.update(chunk))
    readStream.on('end', () => {
      resolve(hash.digest('hex')) // TODO: RO
    })
    readStream.on('error', (error) => {
      reject(error) // TODO: RO
    })
  })
}

// TODO:
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
    // TODO: create string with function?
    let resultPart = `\n${duplicates[i].path}`

    if (duplicates[i].hash !== duplicates[i - 1].hash) {
      resultPart = `\n${resultPart}`
    }

    result += resultPart
  }

  return result
}
