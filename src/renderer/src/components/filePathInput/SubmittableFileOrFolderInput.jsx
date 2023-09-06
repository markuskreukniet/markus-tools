import { createSignal } from 'solid-js'
import ActiveByNumberButton from '../ActiveByNumberButton'
import FileOrFolderInput from './FileOrFolderInput'
import {
  isResultObjectOk,
  toResultObjectWithResultStatusOk
} from '../../../../preload/modules/resultStatus'

export default function SubmittableFileOrFolderInput(props) {
  let selectedFilePathObjects = []
  const [numberOfFilePathObjects, setNumberOfFilePathObjects] = createSignal(0)

  function setState(resultObject) {
    if (isResultObjectOk(resultObject)) {
      selectedFilePathObjects = resultObject.result
      setNumberOfFilePathObjects(selectedFilePathObjects.length)
    } else {
      props.onChange(resultObject)
    }
  }

  function submit() {
    props.onChange(toResultObjectWithResultStatusOk(selectedFilePathObjects))
  }

  const submitButton = (
    <ActiveByNumberButton
      minimumNumber={props.minimumFiles}
      currentNumber={numberOfFilePathObjects()}
      onAction={submit}
      text="submit"
    />
  )

  return <FileOrFolderInput onChange={setState} submitButton={submitButton} />
}
