import { createSignal } from 'solid-js'
import ActivatableSubmitButton from '../activatableButton/ActivatableSubmitButton'
import { Either } from '../../../../preload/monads/either'
import FileOrFolderInput from './FileOrFolderInput'

export default function SubmittableFileOrFolderInput(props) {
  let selectedFileSystemNodes = []
  const [buttonActive, setButtonActive] = createSignal(false)

  function handleChange(result) {
    if (result.isRight()) {
      selectedFileSystemNodes = result.value.selectedFileSystemNodes
      setButtonActive(result.value.hasFileSystemNode)
    } else {
      props.onChange(Either.left(result.value))
    }
  }

  function submit() {
    props.onChange(Either.right(selectedFileSystemNodes))
  }

  const submitButton = <ActivatableSubmitButton active={buttonActive()} onAction={submit} />

  return <FileOrFolderInput onChange={handleChange} submitButton={submitButton} />
}
