import FileOrFolderInput from './FileOrFolderInput'
import { filePathSelectionType } from '../../../../preload/modules/files'
import { isResultObjectOk } from '../../../../preload/modules/resultStatus'

export default function MaxOneDirectoryInput(props) {
  function handleOnChange(resultObject) {
    if (isResultObjectOk(resultObject)) {
      resultObject.result = {
        selectedFilePathObject: resultObject.result.selectedFilePathObjects[0],
        hasFilePathObject: resultObject.result.hasFilePathObject
      }
    }
    props.onChange(resultObject)
  }

  return (
    <FileOrFolderInput
      onChange={handleOnChange}
      filePathSelectionType={filePathSelectionType.directory}
      maxOneInput
    />
  )
}
