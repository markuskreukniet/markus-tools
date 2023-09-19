import { createSignal, For, Show } from 'solid-js'
import ActivatableButton from '../activatableButton/ActivatableButton'
import FilePathSelector from './FilePathSelector'
import { filePathSelectionType } from '../../../../preload/modules/files'
import {
  isResultObjectOk,
  toResultObjectWithResultStatusOk
} from '../../../../preload/modules/resultStatus'

// TODO:
// Adding a file could add a duplicate file since there could already be a folder with its whole tree of child folders already containing that file.
// Adding a folder could add a duplicate file since the folder with its whole tree of child folders could contain a duplicate file.
// Checking child folders of a folder is only possible in the main, which is possible by adding such a function in the main.

export default function FileOrFolderInput(props) {
  const [selectedFilePathObjects, setSelectedFilePathObjects] = createSignal([])
  const [hasFilePathObject, setHasFilePathObject] = createSignal(false)

  function setState(result) {
    if (result.value !== '') {
      if (props.maxOneInput) {
        setSelectedFilePathObjects([result])
      } else if (
        !selectedFilePathObjects().some((filePathObject) => filePathObject.value === result.value)
      ) {
        setSelectedFilePathObjects([...selectedFilePathObjects(), result])
      } else {
        return
      }
      setHasFilePathObject(selectedFilePathObjects().length > 0)
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
    setSelectedFilePathObjects([])
    setHasFilePathObject(false)
  }

  function handleChange(resultObject) {
    if (isResultObjectOk(resultObject)) {
      setState(resultObject.result)
      props.onChange(
        toResultObjectWithResultStatusOk({
          selectedFilePathObjects: selectedFilePathObjects(),
          hasFilePathObject: hasFilePathObject()
        })
      )
    } else {
      props.onChange(resultObject)
    }
  }

  return (
    <div>
      <div class="display-flex not-first-child-margin-left-1">
        <Show when={showFilePathSelector(filePathSelectionType.file)}>
          <FilePathSelector onChange={handleChange} />
        </Show>
        <Show when={showFilePathSelector(filePathSelectionType.directory)}>
          <FilePathSelector onChange={handleChange} directory />
        </Show>
      </div>
      <div class="display-flex justify-content-flex-end not-first-child-margin-left-1">
        <ActivatableButton text="reset" active={hasFilePathObject()} onAction={resetState} />
        {props.submitButton}
      </div>
      <ul>
        <For each={selectedFilePathObjects()}>
          {(filePathObject) => <li>{filePathObject.value}</li>}
        </For>
      </ul>
    </div>
  )
}
