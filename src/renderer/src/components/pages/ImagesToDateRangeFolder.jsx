import { createSignal } from 'solid-js'
import TextResultPage from '../page/TextResultPage'
import ActivatableSubmitButton from '../activatableButton/ActivatableSubmitButton'
import FileOrFolderInput from '../filePathInput/FileOrFolderInput'
import { isResultObjectOk } from '../../../../preload/modules/resultStatus'
import MaxOneDirectoryInput from '../filePathInput/MaxOneDirectoryInput'

export default function imagesToDateRangeFolder(props) {
  let inputFilePathObjects = []
  let outputFilePath = ''
  const [getOutput, setGetOutput] = createSignal(function () {})
  const [hasValidInput, setHasValidInput] = createSignal(false)
  const [status, setStatus] = createSignal('')

  async function setStateWithBE(filePathObjects, path) {
    // TODO: should come from GUI
    const useDirectoriesTreeInput = true

    const imagesToDateRangeFolderRO = await window.dateRangeFolder.imagesToDateRangeFolderBE(
      filePathObjects,
      path,
      useDirectoriesTreeInput
    )
    setStatus(
      isResultObjectOk(imagesToDateRangeFolderRO) ? 'done' : imagesToDateRangeFolderRO.message
    )
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
      setStatus(result.value)
    }
  }

  function handleOutputDirectory(result) {
    if (result.isRight()) {
      outputFilePath = result.value
      validateInput()
    } else {
      setStatus(result.value)
    }
  }

  function submit() {
    setGetOutput(setStateWithBE(inputFilePathObjects, outputFilePath))
  }

  const inputComponent = (
    <div>
      <FileOrFolderInput onChange={handleInputFileSystemNodes} />
      <MaxOneDirectoryInput onChange={handleOutputDirectory} />
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
