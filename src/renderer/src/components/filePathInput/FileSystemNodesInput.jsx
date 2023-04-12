import { createSignal, For, Show } from 'solid-js'
import ActivatableButton from '../activatableButton/ActivatableButton'
import { Either } from '../../../../preload/monads/either'
import FilePathSelector from './FilePathSelector'
import { filePathSelectionType } from '../../../../preload/modules/files'

export default function FileSystemNodesInput(props) {
  const [selectedFileSystemNodes, setSelectedFileSystemNodes] = createSignal([])
  const [hasFileSystemNode, setHasFileSystemNode] = createSignal(false)

  // A trailing slash is needed. Without the slash, /path/sub is a parent of /path/subpath.
  // This trailing slash method should also work on non-Windows systems.
  function getPathWithPossibleTrailingSlash(fileSystemNode) {
    let result = fileSystemNode.path
    if (fileSystemNode.isDirectory) {
      const forwardSlash = '/'
      if (result.startsWith(forwardSlash)) {
        result += forwardSlash
      } else {
        result += '\\'
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

  // We cannot change changedSelectedFileSystemNodes and foundOrDescendantFilePath to one boolean.
  function handleChange(result) {
    if (result.isRight()) {
      if (result.value.path !== '') {
        let changedSelectedFileSystemNodes = false
        if (props.maxOneInput) {
          setSelectedFileSystemNodes([result.value])
          changedSelectedFileSystemNodes = true
        } else {
          const newPath = getPathWithPossibleTrailingSlash(result.value)
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
            setSelectedFileSystemNodes([...filteredSelectedFileSystemNodes, result.value])
            changedSelectedFileSystemNodes = true
          }
        }
        if (changedSelectedFileSystemNodes) {
          setHasFileSystemNode(selectedFileSystemNodes().length > 0)
          onChangeEitherRight()
        }
      }
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
      <div class="file-system-nodes-input__file-path-selector-container">
        <Show when={showFilePathSelector(filePathSelectionType.file)}>
          <FilePathSelector onChange={handleChange} />
        </Show>
        <Show when={showFilePathSelector(filePathSelectionType.directory)}>
          <FilePathSelector onChange={handleChange} directory />
        </Show>
      </div>
      <div class="file-system-nodes-input__submission-buttons">
        <ActivatableButton
          text="reset"
          active={hasFileSystemNode()}
          onAction={resetState}
          variant={'tertiary'}
        />
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
