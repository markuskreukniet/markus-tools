import { createSignal, For, Show } from 'solid-js'
import ActivatableButton from '../activatableButton/ActivatableButton'
import FilePathSelector from './FilePathSelector'
import { filePathSelectionType } from '../../../../preload/modules/files'
import { toResultObjectWithResultStatusOk } from '../../../../preload/modules/resultStatus'

// TODO:
// Adding a file could add a duplicate file since there could already be a folder with its whole tree of child folders already containing that file.
// Adding a folder could add a duplicate file since the folder with its whole tree of child folders could contain a duplicate file.
// Checking child folders of a folder is only possible in the main, which is possible by adding such a function in the main.

export default function FileOrFolderInput(props) {
  const [selectedFileSystemNodes, setSelectedFileSystemNodes] = createSignal([])
  const [hasFileSystemNode, setHasFileSystemNode] = createSignal(false)

  function setState(result) {
    if (result.path !== '') {
      if (props.maxOneInput) {
        setSelectedFileSystemNodes([result])
      } else if (
        !selectedFileSystemNodes().some((fileSystemNode) => fileSystemNode.path === result.path)
      ) {
        setSelectedFileSystemNodes([...selectedFileSystemNodes(), result])
      } else {
        return
      }
      setHasFileSystemNode(selectedFileSystemNodes().length > 0)
    }
  }

  function showFilePathSelector(type) {
    return (
      !props.filePathSelectionType ||
      props.filePathSelectionType === filePathSelectionType.both ||
      props.filePathSelectionType === type
    )
  }

  function resetState() {
    setSelectedFileSystemNodes([])
    setHasFileSystemNode(false)
  }

  // TODO: use Either
  function handleChange(result) {
    if (result.isRight()) {
      setState(result.value)
      props.onChange(
        toResultObjectWithResultStatusOk({
          selectedFileSystemNodes: selectedFileSystemNodes(),
          hasFileSystemNode: hasFileSystemNode()
        })
      )
    } else {
      props.onChange(result)
    }
  }

  return (
    <div>
      <div class="display-flex gap-1">
        <Show when={showFilePathSelector(filePathSelectionType.file)}>
          <FilePathSelector onChange={handleChange} />
        </Show>
        <Show when={showFilePathSelector(filePathSelectionType.directory)}>
          <FilePathSelector onChange={handleChange} directory />
        </Show>
      </div>
      <div class="display-flex justify-content-flex-end gap-1">
        <ActivatableButton text="reset" active={hasFileSystemNode()} onAction={resetState} />
        {props.submitButton}
      </div>
      <ul>
        <For each={selectedFileSystemNodes()}>
          {(FileSystemNode) => <li>{FileSystemNode.path}</li>}
        </For>
      </ul>
    </div>
  )
}
