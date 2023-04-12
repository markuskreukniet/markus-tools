import { createSignal } from 'solid-js'
import TextResultPage from '../page/TextResultPage'
import SubmittableFileSystemNodeInputs from '../filePathInput/SubmittableFileSystemNodeInputs'

import { toEitherRightResult } from '../../../../preload/monads/either'

export default function filesToDateRangeDirectory(props) {
  let inputFilePathObjects = [] // TODO: correct naming?
  let outputFilePath = '' // TODO: correct naming?
  const [eitherResultOutput, setEitherResultOutput] = createSignal(null)
  const [getOutput, setGetOutput] = createSignal(function () {})
  const [hasValidInput, setHasValidInput] = createSignal(false)

  async function setStateWithBE(uniqueFileSystemNodes, destinationDirectoryFilePath) {
    // TODO: result
    await window.goBackend.goFunctionCallBE('filesToDateRangeDirectoryToJSON', {
      uniqueFileSystemNodes,
      destinationDirectoryFilePath
    })

    setEitherResultOutput(toEitherRightResult(null))
  }

  function validateInput() {
    if (inputFilePathObjects.length > 0 && outputFilePath !== '') {
      setHasValidInput(true)
    } else {
      setHasValidInput(false)
    }
  }

  // We could extract similar code of the functions handleInputFileSystemNodes and handleOutputDirectory, for example, to the function handleChange.
  // With this extraction, handleInputFileSystemNodes and handleOutputDirectory call both handleChange.
  // However, this extraction hurts the performance and results in more code.
  function handleInputFileSystemNodes(result) {
    if (result.isRight()) {
      inputFilePathObjects = result.value.selectedFileSystemNodes
      validateInput()
    } else {
      setEitherResultOutput(result)
    }
  }

  function handleOutputDirectory(result) {
    if (result.isRight()) {
      outputFilePath = result.value
      validateInput()
    } else {
      setEitherResultOutput(result)
    }
  }

  function submit() {
    setGetOutput(setStateWithBE(inputFilePathObjects, outputFilePath))
  }

  const fileSystemNodesInputs = [
    { fileSystemNodesInputType: 'all', onChange: handleInputFileSystemNodes },
    { fileSystemNodesInputType: 'maxOneDirectory', onChange: handleOutputDirectory }
  ]

  // TODO: rename filePathInput directory in src\components to FileSystemNodesInput
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
