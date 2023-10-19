import { createSignal } from 'solid-js'
import TextResultPage from '../page/TextResultPage'
import ActivatableSubmitButton from '../activatableButton/ActivatableSubmitButton'
import MaxOneDirectoryInput from '../filePathInput/MaxOneDirectoryInput'

export default function SynchronizeDirectory(props) {
  let originalDirectoryFilePathObject = null
  let destinationDirectoryFilePathObject = null
  const [getOutput, setGetOutput] = createSignal(function () {})
  const [status, setStatus] = createSignal('')

  async function test() {
    const testA = await window.synchronization.synchronizeDirectoryBE(
      originalDirectoryFilePathObject,
      destinationDirectoryFilePathObject
    )
    setStatus(testA)
  }

  function handleInputOriginalDirectoryRO(resultObject) {
    originalDirectoryFilePathObject = resultObject.result.selectedFilePathObject
  }

  function handleInputDestinationDirectoryRO(resultObject) {
    destinationDirectoryFilePathObject = resultObject.result.selectedFilePathObject
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
