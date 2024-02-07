import { objectToJSONAndStringsToGoFunctionCallWithArguments } from './utils/utils'

// TODO: is async needed? same for synchronizeDirectoryTrees
// TODO: use same function as synchronizeDirectoryTrees?
export default async function plainTextFilesToText(uniqueFileSystemNodes) {
  return objectToJSONAndStringsToGoFunctionCallWithArguments('plainTextFilesToTextToJSON', {
    uniqueFileSystemNodes
  })
}
