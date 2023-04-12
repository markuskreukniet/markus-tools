import { createSignal } from 'solid-js'
import TextResultPage from '../page/TextResultPage'
import SubmittableFileSystemNodeInputs from '../filePathInput/SubmittableFileSystemNodeInputs'

export default function SynchronizeDirectoryTrees(props) {
  let sourceDirectoryFilePath = ''
  let destinationDirectoryFilePath = ''
  // TODO: should be '' instead of null? Same possible problem in other files
  const [eitherResultOutput, setEitherResultOutput] = createSignal(null)
  const [getOutput, setGetOutput] = createSignal(function () {})
  const [hasValidInput, setHasValidInput] = createSignal(false)

  async function setStateWithBE() {
    const result = await window.goBackend.goFunctionCallBE('synchronizeDirectoryTreesToJSON', {
      sourceDirectoryFilePath,
      destinationDirectoryFilePath
    })
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

  const fileSystemNodesInputs = [
    { fileSystemNodesInputType: 'maxOneDirectory', onChange: handleInputSourceDirectory },
    { fileSystemNodesInputType: 'maxOneDirectory', onChange: handleInputDestinationDirectory }
  ]

  const inputComponent = (
    <SubmittableFileSystemNodeInputs
      fileSystemNodesInputs={fileSystemNodesInputs}
      hasValidInput={hasValidInput()}
      onAction={submit}
    />
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
