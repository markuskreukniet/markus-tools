import { objectToJSONAndStringsToGoFunctionCallWithArguments } from './utils/utils'

export default async function synchronizeDirectoryTrees(sourceDirectory, destinationDirectory) {
  return objectToJSONAndStringsToGoFunctionCallWithArguments('synchronizeDirectoryTreesToJSON', {
    sourceDirectory,
    destinationDirectory
  })
}
