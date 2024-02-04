import { createSignal } from 'solid-js'
import TextResultPage from '../page/TextResultPage'
import ActivatableSubmitButton from '../activatableButton/ActivatableSubmitButton'
import FileSystemNodesInput from '../filePathInput/FileSystemNodesInput'
import MaxOneDirectoryInput from '../filePathInput/MaxOneDirectoryInput'
import { toEitherRightResult } from '../../../../preload/monads/either'

export default function imagesToDateRangeFolder(props) {
  let inputFilePathObjects = []
  let outputFilePath = ''
  const [eitherResultOutput, setEitherResultOutput] = createSignal(null)
  const [getOutput, setGetOutput] = createSignal(function () {})
  const [hasValidInput, setHasValidInput] = createSignal(false)

  async function setStateWithBE(filePathObjects, filePath) {
    const imagesToDateRangeFolderRO = await window.dateRangeFolder.imagesToDateRangeFolderBE(
      filePathObjects,
      filePath
    )
    // TODO: imagesToDateRangeFolderBE should return eitherResult and should be setEitherResultOutput(result)
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

  const inputComponent = (
    <div>
      <FileSystemNodesInput onChange={handleInputFileSystemNodes} />
      <MaxOneDirectoryInput onChange={handleOutputDirectory} />
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
