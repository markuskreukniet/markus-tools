import { createSignal, For, Show } from 'solid-js'
import ActiveByNumberButton from './ActiveByNumberButton'
import FilePathSelector from './filePathInput/FilePathSelector'
import { filePathSelectionType } from '../../../preload/modules/files'
import {
  isResultObjectOk,
  toResultObjectWithResultStatusOk
} from '../../../preload/modules/resultStatus'

// TODO: also files to folder FilePathInput in components
export default function SubmittableFileOrFolderInput(props) {
  const [selectedFilePathObjects, setSelectedFilePathObjects] = createSignal([])
  const [numberOfFilePathObjects, setNumberOfFilePathObjects] = createSignal(0)

  function setState(resultObject) {
    if (isResultObjectOk(resultObject)) {
      if (resultObject.result.value !== '') {
        if (props.maxOneInput) {
          setSelectedFilePathObjects([resultObject.result])
        } else if (
          !selectedFilePathObjects().some(
            (filePathObject) => filePathObject.value === resultObject.result.value
          )
        ) {
          setSelectedFilePathObjects([...selectedFilePathObjects(), resultObject.result])
        } else {
          return
        }
        setNumberOfFilePathObjects(selectedFilePathObjects().length)
      }
    } else {
      props.onChange(resultObject)
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
    setNumberOfFilePathObjects(0)
  }

  function submit() {
    props.onChange(toResultObjectWithResultStatusOk(selectedFilePathObjects()))
  }

  return (
    <div>
      <div class="display-flex not-first-child-margin-left-1">
        <Show when={showFilePathSelector(filePathSelectionType.file)}>
          <FilePathSelector onChange={setState} />
        </Show>
        <Show when={showFilePathSelector(filePathSelectionType.directory)}>
          <FilePathSelector onChange={setState} directory />
        </Show>
      </div>
      <div class="display-flex justify-content-flex-end not-first-child-margin-left-1">
        <ActiveByNumberButton
          minimumNumber={1}
          currentNumber={numberOfFilePathObjects()}
          onAction={resetState}
          text="reset"
        />
        <ActiveByNumberButton
          minimumNumber={props.minimumFiles}
          currentNumber={numberOfFilePathObjects()}
          onAction={submit}
          text="submit"
        />
      </div>
      <ul>
        <For each={selectedFilePathObjects()}>
          {(filePathObject) => <li>{filePathObject.value}</li>}
        </For>
      </ul>
    </div>
  )
}
