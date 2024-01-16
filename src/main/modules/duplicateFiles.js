import crypto from 'crypto'
import { filePathObjectsToFileObjects, getReadFileHandle } from './filePaths.js'
import { concatErrorMessageOnNewLine } from '../../preload/modules/errors'
import {
  isResultObjectOk,
  toResultObjectWithNullResultAndResultStatusErrorSystem,
  toResultObjectWithResultStatusOk
} from '../../preload/modules/resultStatus'

// TODO: this function (and others) does not have to async
export default async function duplicateFiles(filePathObjects) {
  const filePathObjectsToFileObjectsRO = await filePathObjectsToFileObjects(filePathObjects, true)
  if (!isResultObjectOk(filePathObjectsToFileObjectsRO)) {
    return filePathObjectsToFileObjectsRO
  }

  const fileObjects = filePathObjectsToFileObjectsRO.result
  fileObjects.sort(compare)

  const duplicates = []
  let lastPushedIndex = -1

  for (let i = 1; i < fileObjects.length; i++) {
    const fileObject = fileObjects[i - 1]
    const fileObject2 = fileObjects[i]

    if (fileObject.size === fileObject2.size) {
      // Performing two times getFileHash for the same file is a performance loss, which happens otherwise with at least three files with the same file size.
      const fileHash = fileObject.fileHash || (await getFileHash(fileObject.path))
      const fileHash2 = await getFileHash(fileObject2.path)
      fileObject2.fileHash = fileHash2

      if (fileHash === fileHash2) {
        if (lastPushedIndex !== i - 1) {
          duplicates.push(toDuplicateObject(fileObject.path, fileHash))
        }
        duplicates.push(toDuplicateObject(fileObject2.path, fileHash2))
        lastPushedIndex = i
      }
    }
  }

  if (duplicates.length === 0) {
    return ''
  } else {
    return duplicatesArrayToResultString(duplicates)
  }
}

function toDuplicateObject(path, hash) {
  return { path, hash }
}

async function getFileHash(filePath) {
  const fileHandleRO = await getReadFileHandle(filePath)
  if (!isResultObjectOk(fileHandleRO)) {
    return fileHandleRO
  }

  let readStream = null
  let resultRO = null
  try {
    readStream = fileHandleRO.result.createReadStream()
  } catch (error) {
    resultRO = toResultObjectWithNullResultAndResultStatusErrorSystem(error.message)
  }

  if (!resultRO) {
    resultRO = await getFileHashByReadStream(readStream)

    // When we use "readStream.on 'error'," we don't have to use try-catch
    // readStream.destroy() should come before filehandle.close()
    readStream.destroy()
  }

  try {
    await fileHandleRO.result.close()
  } catch (error) {
    if (isResultObjectOk(resultRO)) {
      return toResultObjectWithNullResultAndResultStatusErrorSystem(error.message)
    } else {
      resultRO.message = concatErrorMessageOnNewLine(resultRO.message, error.message)
    }
  }

  return resultRO.result // TODO:
}

function getFileHashByReadStream(readStream) {
  // sha256 is generally faster and more secure than SHA1
  // SHA1 is generally faster and more secure than MD5
  const hash = crypto.createHash('sha256')

  return new Promise((resolve, reject) => {
    readStream.on('data', (chunk) => hash.update(chunk))
    readStream.on('end', () => {
      resolve(toResultObjectWithResultStatusOk(hash.digest('hex')))
    })
    readStream.on('error', (error) => {
      reject(toResultObjectWithNullResultAndResultStatusErrorSystem(error.message))
    })
  })
}

function compare(a, b) {
  if (a.size < b.size) {
    return -1
  }
  if (a.size > b.size) {
    return 1
  }
  return 0
}

function prependNewline(string) {
  return `\n${string}`
}

function duplicatesArrayToResultString(duplicates) {
  let result = duplicates[0].path
  for (let i = 1; i < duplicates.length; i++) {
    let resultPart = prependNewline(duplicates[i].path)
    if (duplicates[i].hash !== duplicates[i - 1].hash) {
      resultPart = prependNewline(resultPart)
    }
    result += resultPart
  }

  return result
}
