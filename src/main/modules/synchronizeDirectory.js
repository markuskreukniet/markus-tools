import {
  combineOutputFilePathWithRelativeInputFilePath,
  copyFile,
  filePathExists,
  getFileAndDirectoryFileObjects,
  makeDirectory,
  removeDirectoryTree,
  removeFile
} from './filePaths.js'
import { isResultObjectOk } from '../../preload/modules/resultStatus'

export default async function synchronizeDirectory(
  originalDirectoryFilePath,
  destinationDirectoryFilePath
) {
  const directoriesTree = true

  const stack = [originalDirectoryFilePath]
  while (stack.length > 0) {
    const originalDirectoryPath = stack.pop()
    const destinationDirectoryPath = originalDirectoryPath.replace(
      originalDirectoryFilePath,
      destinationDirectoryFilePath
    )

    const originalFileAndDirectoryFileObjectsRO = await getFileAndDirectoryFileObjects(
      originalDirectoryPath,
      directoriesTree
    )
    if (!isResultObjectOk(originalFileAndDirectoryFileObjectsRO)) {
      return originalFileAndDirectoryFileObjectsRO
    }

    // TODO: originalFileAndDirectoryFileObjectsRO.result should be already a Map
    // const originalFileAndDirectoryFileObjectsPathMap = new Map(
    //   originalFileAndDirectoryFileObjectsRO.result.map((fileObject) => [
    //     fileObject.path,
    //     {
    //       dateModified: fileObject.dateCreated,
    //       isDirectory: fileObject.isDirectory
    //     }
    //   ])
    // )

    // TODO: almost the same as originalFileAndDirectoryFileObjectsRO
    const destinationFileAndDirectoryFileObjectsRO = await getFileAndDirectoryFileObjects(
      destinationDirectoryPath,
      directoriesTree
    )
    if (!isResultObjectOk(destinationFileAndDirectoryFileObjectsRO)) {
      return destinationFileAndDirectoryFileObjectsRO
    }

    // TODO: should be already a Map
    const destinationFileAndDirectoryFileObjectsPathMap = new Map(
      destinationFileAndDirectoryFileObjectsRO.result.map((fileObject) => [
        fileObject.path,
        {
          dateModified: fileObject.dateCreated,
          isDirectory: fileObject.isDirectory
        }
      ])
    )

    for (const originalFileObject of originalFileAndDirectoryFileObjectsRO.result) {
      const destinationFilePath = combineOutputFilePathWithRelativeInputFilePath(
        originalDirectoryFilePath,
        originalFileObject.path,
        destinationDirectoryFilePath
      )

      if (destinationFileAndDirectoryFileObjectsPathMap.has(destinationFilePath)) {
        if (
          !originalFileObject.isDirectory &&
          originalFileObject.dateCreated >
            destinationFileAndDirectoryFileObjectsPathMap.get(destinationFilePath).dateModified
        ) {
          // TODO: RO
          // copyFile does replace, fs.copyFile and fs.createWriteStream, both do that, keep this comment, but in filePaths.js
          await copyFile(originalFileObject.path, destinationFilePath)
        }

        destinationFileAndDirectoryFileObjectsPathMap.delete(destinationFilePath)
      } else {
        if (originalFileObject.isDirectory) {
          // TODO: RO
          await makeDirectory(destinationFilePath)
        } else {
          // TODO: RO and copied
          await copyFile(originalFileObject.path, destinationFilePath)
        }
      }
    }

    for (const [
      destinationFilePath,
      destinationFileObject
    ] of destinationFileAndDirectoryFileObjectsPathMap.entries()) {
      const originalFilePath = combineOutputFilePathWithRelativeInputFilePath(
        destinationDirectoryFilePath,
        destinationFilePath,
        originalDirectoryFilePath
      )

      const filePathExistsRO = await filePathExists(originalFilePath)
      if (!isResultObjectOk(filePathExistsRO)) {
        return filePathExistsRO
      }

      if (!filePathExistsRO.result) {
        if (destinationFileObject.isDirectory) {
          // TODO: RO
          await removeDirectoryTree(destinationFilePath)
        } else {
          // TODO: RO
          await removeFile(destinationFilePath)
        }
      }
    }
  }

  return `${originalDirectoryFilePath} testB ${destinationDirectoryFilePath}`
}

// Go test code
// const { exec } = require('child_process')

//   const jsonArguments = JSON.stringify({
//     sourceDirectoryFilePath,
//     destinationDirectoryFilePath
//   }).replace(/"/g, '\\"')
//   const goProcess = exec(`go run ./go/main.go "${jsonArguments}"`, (error, stdout) => {
//     if (error) {
//       console.error(`Error executing Go program: ${error}`)
//       return
//     }
//     console.log(`Go program output: ${stdout}`)
//   })

//   goProcess.on('close', (code) => {
//     console.log(`Go program exited with code ${code}`)
//   })
