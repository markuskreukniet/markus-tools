import { createSignal, For } from 'solid-js'
import ActiveByNumberButton from './ActiveByNumberButton'
import FilePathSelector from './FilePathSelector'
import {
  isResultObjectOk,
  toResultObjectWithResultStatusOk
} from '../../../preload/modules/resultStatus'

// TODO:
// Adding a file could add a duplicate file since there could already be a folder with its whole tree of child folders already containing that file.
// Adding a folder could add a duplicate file since the folder with its whole tree of child folders could contain a duplicate file.

// Checking child folders of a folder is only possible in the main, which is possible by adding such a function in the main.

export default function FileOrFolderInput(props) {
  const [selectedFilePathObjects, setSelectedFilePathObjects] = createSignal([])
  const [numberOfFilePathObjects, setNumberOfFilePathObjects] = createSignal(0)

  function setState(resultObject) {
    if (isResultObjectOk(resultObject)) {
      if (
        resultObject.result.value !== '' &&
        !selectedFilePathObjects().some(
          (filePathObject) => filePathObject.value === resultObject.result.value
        )
      ) {
        setSelectedFilePathObjects([...selectedFilePathObjects(), resultObject.result])
        setNumberOfFilePathObjects(selectedFilePathObjects().length)
      }
    } else {
      props.onChange(resultObject)
    }
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
        <FilePathSelector onChange={setState} />
        <FilePathSelector onChange={setState} directory />
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
