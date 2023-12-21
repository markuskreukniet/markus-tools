import { createSignal, For, Show } from 'solid-js'
import ActivatableButton from '../activatableButton/ActivatableButton'
import { Either } from '../../../../preload/monads/either'
import FilePathSelector from './FilePathSelector'
import { filePathSelectionType } from '../../../../preload/modules/files'

// TODO:
// Adding a file could add a duplicate file since there could already be a directory with its whole tree of child directories already containing that file.
// Adding a directory could add a duplicate file since the directory with its whole tree of child directories could contain a duplicate file.
// Checking child directories of a directory is only possible in the main, which is possible by adding such a function in the main.

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

  // C:/development/markus-tools
  // C:/development/markus-tools/go
  // C:/development/markus-tools/test.go

  // in:  C:/development/markus-tools
  // new: C:/development/markus-tools/test.go

  // in:  C:/development/markus-tools/test.go, C:/development/test
  // new: C:/development/markus-tools

  // TODO: name and use function, or function content
  // A trailing slash is needed. Without the slash, /path/sub is a parent of /path/subpath.
  // This trailing slash method should also work on non-Windows systems.
  function asdf(newFileSystemNode) {
    const slash = '/'
    const newPath = newFileSystemNode.path + slash
    const filteredSelectedFileSystemNodes = []
    let foundOrDescendantFilePath = false
    for (const node of selectedFileSystemNodes()) {
      const nodePath = node.path + slash
      if (newPath === nodePath || newPath.startsWith(nodePath)) {
        foundOrDescendantFilePath = true
        break
      }
      if (!nodePath.startsWith(newPath)) {
        filteredSelectedFileSystemNodes.push(node)
      }
    }
    if (!foundOrDescendantFilePath) {
      setSelectedFileSystemNodes([...filteredSelectedFileSystemNodes, newFileSystemNode])
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
