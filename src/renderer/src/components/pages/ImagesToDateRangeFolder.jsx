import { createSignal } from 'solid-js'
import TextResultPage from '../page/TextResultPage'
import ActivatableSubmitButton from '../ActivatableSubmitButton'
import FileOrFolderInput from '../filePathInput/FileOrFolderInput'
import { filePathSelectionType } from '../../../../preload/modules/files'
import { isResultObjectOk } from '../../../../preload/modules/resultStatus'

export default function imagesToDateRangeFolder(props) {
  let inputFilePathObjects = []
  let outputFilePath = ''
  const [getOutput, setGetOutput] = createSignal(function () {})
  const [hasValidInput, setHasValidInput] = createSignal(false)
  const [status, setStatus] = createSignal('')

  async function processInputToOutput(filePathObjects, path) {
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

  // We could extract similar code of the functions handleInputFilePathsRO and handleOutputDirectoryRO, for example, to the function handleRO.
  // With this extraction, handleInputFilePathsRO and handleOutputDirectoryRO call both handleRO.
  // However, this extraction hurts the performance and results in more code.
  function handleInputFilePathsRO(resultObject) {
    if (isResultObjectOk(resultObject)) {
      inputFilePathObjects = resultObject.result
      validateInput()
    } else {
      setStatus(resultObject.message)
    }
  }

  function handleOutputDirectoryRO(resultObject) {
    if (isResultObjectOk(resultObject)) {
      // TODO: should not return an array when maxOneInput?
      outputFilePath = resultObject.result[0].value
      validateInput()
    } else {
      setStatus(resultObject.message)
    }
  }

  function submit() {
    setGetOutput(processInputToOutput(inputFilePathObjects, outputFilePath))
  }

  const inputComponent = (
    <div>
      <FileOrFolderInput onChange={handleInputFilePathsRO} />
      <FileOrFolderInput
        onChange={handleOutputDirectoryRO}
        filePathSelectionType={filePathSelectionType.directory}
        maxOneInput
      />
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
