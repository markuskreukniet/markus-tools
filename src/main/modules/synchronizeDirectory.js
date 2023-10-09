import { copyDirectoryTree, copyFile, getFileAndDirectoryFileObjects } from './filePaths.js'

export default async function synchronizeDirectory(
  originalDirectoryFilePathObject,
  destinationDirectoryFilePathObject
) {
  // TODO: this boolean should come from UI
  const directoriesTree = true
  // added getFileAndDirectoryFileObjects for synchronizeDirectory

  const stack = [originalDirectoryFilePathObject.value]
  while (stack.length > 0) {
    const originalDirectoryPath = stack.pop()
    const destinationDirectoryPath = originalDirectoryPath.replace(
      originalDirectoryFilePathObject.value,
      destinationDirectoryFilePathObject.value
    )

    const originalFileAndDirectoryFileObjectsRO = await getFileAndDirectoryFileObjects(
      originalDirectoryPath,
      directoriesTree
    )
    const destinationFileAndDirectoryFileObjectsRO = await getFileAndDirectoryFileObjects(
      destinationDirectoryPath,
      directoriesTree
    )

    for (const fileObject of originalFileAndDirectoryFileObjectsRO.result) {
      if (fileObject.isDirectory) {
        // if destination does not have the directory, copyDirectoryTree
        await copyDirectoryTree()
        // else stack.push(fileObject.path)
      } else {
        // if destination does not have the file
        await copyFile()
        // if destination does have the file and original file is newer, replace the file
        // TODO: copyFile does replace?
        await copyFile()
      }
    }
  }
}
