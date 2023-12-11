import { createSignal } from 'solid-js'
import TextResultPage from '../page/TextResultPage'
import ActivatableSubmitButton from '../activatableButton/ActivatableSubmitButton'
import {
  eitherLeftResultToErrorString,
  isEitherRightResult
} from '../../../../preload/monads/either'
import MaxOneDirectoryInput from '../filePathInput/MaxOneDirectoryInput'

export default function SynchronizeDirectory(props) {
  let originalDirectoryFilePath = null
  let destinationDirectoryFilePath = null
  const [getOutput, setGetOutput] = createSignal(function () {})
  const [status, setStatus] = createSignal('')

  // TODO: rename synchronizeDirectory to what the Go version is
  // TODO: rename also on other to callBE()
  async function callBE() {
    const result = await window.synchronization.synchronizeDirectoryBE(
      originalDirectoryFilePath,
      destinationDirectoryFilePath
    )
    if (isEitherRightResult(result)) {
      // TODO: done is also used somewhere else
      setStatus('done')
    } else {
      setStatus(eitherLeftResultToErrorString(result))
    }
  }

  // TODO: either
  function handleInputSourceDirectory(result) {
    originalDirectoryFilePath = result.value
  }

  function handleInputDestinationDirectory(result) {
    destinationDirectoryFilePath = result.value
  }

  function submit() {
    setGetOutput(callBE)
  }

  const inputComponent = (
    <div>
      <MaxOneDirectoryInput onChange={handleInputSourceDirectory} />
      <MaxOneDirectoryInput onChange={handleInputDestinationDirectory} />
      <ActivatableSubmitButton active={true} onAction={submit} />
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
