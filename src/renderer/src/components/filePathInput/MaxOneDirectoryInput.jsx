import FileOrFolderInput from './FileOrFolderInput'
import { filePathSelectionType } from '../../../../preload/modules/files'
import { isResultObjectOk } from '../../../../preload/modules/resultStatus'

export default function MaxOneDirectoryInput(props) {
  function handleChange(resultObject) {
    if (isResultObjectOk(resultObject)) {
      resultObject.result = {
        selectedFileSystemNode: resultObject.result.selectedFileSystemNodes[0],
        hasFileSystemNode: resultObject.result.hasFileSystemNode
      }
    }
    props.onChange(resultObject)
  }

  return (
    <FileOrFolderInput
      onChange={handleChange}
      filePathSelectionType={filePathSelectionType.directory}
      maxOneInput
    />
  )
}
