import { getFileAndDirectoryFileObjects } from './filePaths.js'

export default async function synchronizeDirectory(
  originalDirectoryFilePathObject,
  destinationDirectoryFilePathObject
) {
  // TODO: this boolean should come from UI
  const directoriesTree = true

  // added getFileAndDirectoryFileObjects for synchronizeDirectory
  getFileAndDirectoryFileObjects(originalDirectoryFilePathObject.value, directoriesTree)

  getFileAndDirectoryFileObjects(destinationDirectoryFilePathObject, directoriesTree)

  // if a file in the destination directory exists and if the date modified of the file in the original directory is newer,
  // then replace the file in the destination directory
}
