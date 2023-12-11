import FileOrFolderInput from './FileOrFolderInput'
import { filePathSelectionType } from '../../../../preload/modules/files'

export default function MaxOneDirectoryInput(props) {
  function handleChange(result) {
    if (result.isRight()) {
      // result.value.hasFileSystemNode is not needed here at the moment
      if (result.value.selectedFileSystemNodes.length === 1) {
        result.value = result.value.selectedFileSystemNodes[0].path
      } else {
        result.value = ''
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
