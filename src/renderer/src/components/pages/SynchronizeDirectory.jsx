import { createSignal } from 'solid-js'
import TextResultPage from '../page/TextResultPage'
import ActivatableSubmitButton from '../activatableButton/ActivatableSubmitButton'
import MaxOneDirectoryInput from '../filePathInput/MaxOneDirectoryInput'

export default function SynchronizeDirectory(props) {
  let originalDirectoryFilePath = null
  let destinationDirectoryFilePath = null
  const [getOutput, setGetOutput] = createSignal(function () {})
  const [status, setStatus] = createSignal('')

  async function test() {
    const testA = await window.synchronization.synchronizeDirectoryBE(
      originalDirectoryFilePath,
      destinationDirectoryFilePath
    )
    setStatus(testA)
  }

  function handleInputOriginalDirectoryRO(resultObject) {
    // TODO: resultObject.result should be different, maybe only value?
    originalDirectoryFilePath = resultObject.result.selectedFilePathObject.value
  }

  function handleInputDestinationDirectoryRO(resultObject) {
    destinationDirectoryFilePath = resultObject.result.selectedFilePathObject.value
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
