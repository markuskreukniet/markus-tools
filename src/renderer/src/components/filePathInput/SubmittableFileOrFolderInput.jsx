import { createSignal } from 'solid-js'
import ActivatableSubmitButton from '../activatableButton/ActivatableSubmitButton'
import FileOrFolderInput from './FileOrFolderInput'
import {
  isResultObjectOk,
  toResultObjectWithResultStatusOk
} from '../../../../preload/modules/resultStatus'

export default function SubmittableFileOrFolderInput(props) {
  let selectedFilePathObjects = []
  const [buttonActive, setButtonActive] = createSignal(false)

  function setState(resultObject) {
    if (isResultObjectOk(resultObject)) {
      selectedFilePathObjects = resultObject.result.selectedFilePathObjects
      setButtonActive(resultObject.result.hasFilePathObject)
    } else {
      props.onChange(resultObject)
    }
  }

  function submit() {
    props.onChange(toResultObjectWithResultStatusOk(selectedFilePathObjects))
  }

  const submitButton = <ActivatableSubmitButton active={buttonActive()} onAction={submit} />

  return <FileOrFolderInput onChange={setState} submitButton={submitButton} />
}
