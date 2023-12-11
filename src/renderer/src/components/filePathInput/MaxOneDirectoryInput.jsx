import FileOrFolderInput from './FileOrFolderInput'
import { filePathSelectionType } from '../../../../preload/modules/files'

export default function MaxOneDirectoryInput(props) {
  function handleChange(result) {
    if (result.isRight()) {
      // result.value.hasFileSystemNode is not needed here at the moment
      result.value = result.value.selectedFileSystemNodes[0].path
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
