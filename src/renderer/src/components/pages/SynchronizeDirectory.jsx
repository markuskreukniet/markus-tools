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
  // TODO: rename test
  async function test() {
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

  function handleInputOriginalDirectoryRO(resultObject) {
    // TODO: resultObject.result should be different, maybe only path?
    originalDirectoryFilePath = resultObject.result.selectedFileSystemNode.path
  }

  function handleInputDestinationDirectoryRO(resultObject) {
    destinationDirectoryFilePath = resultObject.result.selectedFileSystemNode.path
  }

  function submit() {
    setGetOutput(test)
  }

  const inputComponent = (
    <div>
      <MaxOneDirectoryInput onChange={handleInputOriginalDirectoryRO} />
      <MaxOneDirectoryInput onChange={handleInputDestinationDirectoryRO} />
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
