import FileOrFolderInput from './FileOrFolderInput'
import { filePathSelectionType } from '../../../../preload/modules/files'

export default function MaxOneDirectoryInput(props) {
  function handleChange(result) {
    if (result.isRight()) {
      result.value = {
        selectedFileSystemNode: result.value.selectedFileSystemNodes[0],
        hasFileSystemNode: result.value.hasFileSystemNode
      }
    }
    props.onChange(result)
  }

  return (
    <FileOrFolderInput
      onChange={handleChange}
      filePathSelectionType={filePathSelectionType.directory}
      maxOneInput
    />
  )
}
