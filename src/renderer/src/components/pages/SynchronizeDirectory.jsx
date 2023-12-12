import { createSignal } from 'solid-js'
import TextResultPage from '../page/TextResultPage'
import ActivatableSubmitButton from '../activatableButton/ActivatableSubmitButton'
import {
  eitherLeftResultToErrorString,
  isEitherRightResult
} from '../../../../preload/monads/either'
import MaxOneDirectoryInput from '../filePathInput/MaxOneDirectoryInput'

export default function SynchronizeDirectory(props) {
  let sourceDirectoryFilePath = ''
  let destinationDirectoryFilePath = ''
  const [getOutput, setGetOutput] = createSignal(function () {})
  const [hasValidInput, setHasValidInput] = createSignal(false)
  const [status, setStatus] = createSignal('')

  // TODO: rename synchronizeDirectory (also the file and app tab and import synchronizeDirectory) to what the Go version is
  async function setStateWithBE() {
    const result = await window.synchronization.synchronizeDirectoryTreesBE(
      sourceDirectoryFilePath,
      destinationDirectoryFilePath
    )
    if (isEitherRightResult(result)) {
      // TODO: done is also used somewhere else
      setStatus('done')
    } else {
      setStatus(eitherLeftResultToErrorString(result))
    }
  }

  function validateInput() {
    if (sourceDirectoryFilePath && destinationDirectoryFilePath) {
      setHasValidInput(true)
    } else {
      setHasValidInput(false)
    }
  }

  function handleInputSourceDirectory(result) {
    if (result.isRight()) {
      sourceDirectoryFilePath = result.value
      validateInput()
    } else {
      setStatus(eitherLeftResultToErrorString(result))
    }
  }

  function handleInputDestinationDirectory(result) {
    if (result.isRight()) {
      destinationDirectoryFilePath = result.value
      validateInput()
    } else {
      setStatus(eitherLeftResultToErrorString(result))
    }
  }

  function submit() {
    setGetOutput(setStateWithBE)
  }

  const inputComponent = (
    <div>
      <MaxOneDirectoryInput onChange={handleInputSourceDirectory} />
      <MaxOneDirectoryInput onChange={handleInputDestinationDirectory} />
      <ActivatableSubmitButton active={hasValidInput()} onAction={submit} />
    </div>
  )

  return (
    <TextResultPage
      title={props.title}
      inputComponent={inputComponent}
      output={status()}
      getOutput={getOutput}
      onLoading={props.onLoading}
    />
  )
}
