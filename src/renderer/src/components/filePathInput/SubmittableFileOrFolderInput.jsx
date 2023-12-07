import { createSignal } from 'solid-js'
import ActivatableSubmitButton from '../activatableButton/ActivatableSubmitButton'
import FileOrFolderInput from './FileOrFolderInput'
import {
  isResultObjectOk,
  toResultObjectWithResultStatusOk
} from '../../../../preload/modules/resultStatus'

export default function SubmittableFileOrFolderInput(props) {
  let selectedFileSystemNodes = []
  const [buttonActive, setButtonActive] = createSignal(false)

  function setState(resultObject) {
    if (isResultObjectOk(resultObject)) {
      selectedFileSystemNodes = resultObject.result.selectedFileSystemNodes
      setButtonActive(resultObject.result.hasFileSystemNode)
    } else {
      props.onChange(resultObject)
    }
  }

  function submit() {
    props.onChange(toResultObjectWithResultStatusOk(selectedFileSystemNodes))
  }

  const submitButton = <ActivatableSubmitButton active={buttonActive()} onAction={submit} />

  return <FileOrFolderInput onChange={setState} submitButton={submitButton} />
}
