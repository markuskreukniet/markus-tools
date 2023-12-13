import { createSignal } from 'solid-js'
import TextResultPage from '../page/TextResultPage'
import ActivatableSubmitButton from '../activatableButton/ActivatableSubmitButton'
import MaxOneDirectoryInput from '../filePathInput/MaxOneDirectoryInput'

export default function SynchronizeDirectoryTrees(props) {
  let sourceDirectoryFilePath = ''
  let destinationDirectoryFilePath = ''
  const [eitherResultOutput, setEitherResultOutput] = createSignal(null)
  const [getOutput, setGetOutput] = createSignal(function () {})
  const [hasValidInput, setHasValidInput] = createSignal(false)

  async function setStateWithBE() {
    const result = await window.synchronization.synchronizeDirectoryTreesBE(
      sourceDirectoryFilePath,
      destinationDirectoryFilePath
    )
    setEitherResultOutput(result)
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
      setEitherResultOutput(result)
    }
  }

  function handleInputDestinationDirectory(result) {
    if (result.isRight()) {
      destinationDirectoryFilePath = result.value
      validateInput()
    } else {
      setEitherResultOutput(result)
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
      eitherResultOutput={eitherResultOutput()}
      getOutput={getOutput}
      onLoading={props.onLoading}
    />
  )
}
