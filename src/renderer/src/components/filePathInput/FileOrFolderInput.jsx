import { createSignal, For, Show } from 'solid-js'
import ActivatableButton from '../activatableButton/ActivatableButton'
import { Either } from '../../../../preload/monads/either'
import FilePathSelector from './FilePathSelector'
import { filePathSelectionType } from '../../../../preload/modules/files'

export default function FileOrFolderInput(props) {
  const [selectedFileSystemNodes, setSelectedFileSystemNodes] = createSignal([])
  const [hasFileSystemNode, setHasFileSystemNode] = createSignal(false)

  // TODO: changedSelectedFileSystemNodes and foundOrDescendantFilePath could be changed to one bool?
  // A trailing slash is needed. Without the slash, /path/sub is a parent of /path/subpath.
  // This trailing slash method should also work on non-Windows systems.
  function setState(result) {
    if (result.path !== '') {
      let changedSelectedFileSystemNodes = false
      if (props.maxOneInput) {
        setSelectedFileSystemNodes([result])
        changedSelectedFileSystemNodes = true
      } else {
        const newPath = getPathWithPossibleTrailingSlash(result)
        const filteredSelectedFileSystemNodes = []
        let foundOrDescendantFilePath = false
        for (const node of selectedFileSystemNodes()) {
          const nodePath = getPathWithPossibleTrailingSlash(node)
          if (newPath === nodePath || newPath.startsWith(nodePath)) {
            foundOrDescendantFilePath = true
            break
          }
          if (!nodePath.startsWith(newPath)) {
            filteredSelectedFileSystemNodes.push(node)
          }
        }
        if (!foundOrDescendantFilePath) {
          setSelectedFileSystemNodes([...filteredSelectedFileSystemNodes, result])
          changedSelectedFileSystemNodes = true
        }
      }
      if (changedSelectedFileSystemNodes) {
        setHasFileSystemNode(selectedFileSystemNodes().length > 0)
        // TODO: onChange
      }
    }
  }

  function getPathWithPossibleTrailingSlash(fileSystemNode) {
    let result = fileSystemNode.path
    if (fileSystemNode.isDirectory) {
      if (result.startsWith('/')) {
        result = result + '/'
      } else {
        result = result + '\\'
      }
    }
    return result
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
    onChangeEitherRight()
  }

  function handleChange(result) {
    if (result.isRight()) {
      setState(result.value)
      onChangeEitherRight()
    } else {
      props.onChange(result)
    }
  }

  function onChangeEitherRight() {
    props.onChange(
      Either.right({
        selectedFileSystemNodes: selectedFileSystemNodes(),
        hasFileSystemNode: hasFileSystemNode()
      })
    )
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
