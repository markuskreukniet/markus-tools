import { createSignal } from 'solid-js'
import TextResultPage from '../page/TextResultPage'
import ActivatableSubmitButton from '../activatableButton/ActivatableSubmitButton'
import { isEitherRightResult } from '../../../../preload/monads/either'
import MaxOneDirectoryInput from '../filePathInput/MaxOneDirectoryInput'

export default function SynchronizeDirectory(props) {
  let originalDirectoryFilePath = null
  let destinationDirectoryFilePath = null
  const [getOutput, setGetOutput] = createSignal(function () {})
  const [status, setStatus] = createSignal('')

  // TODO: rename test
  async function test() {
    const result = await window.synchronization.synchronizeDirectoryBE(
      originalDirectoryFilePath,
      destinationDirectoryFilePath
    )
    console.log('result', result)
    if (isEitherRightResult(result)) {
      // TODO: done is also used somewhere else
      setStatus('done')
    } else {
      // TODO: use function that returns error string
      setStatus('')
    }
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
