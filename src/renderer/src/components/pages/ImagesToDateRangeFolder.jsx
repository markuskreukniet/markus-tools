import { createSignal } from 'solid-js'
import TextResultPage from '../page/TextResultPage'
import ActiveByNumberButton from '../ActiveByNumberButton'
import FileOrFolderInput from '../filePathInput/FileOrFolderInput'
import { filePathSelectionType } from '../../../../preload/modules/files'
import SubmittableFileOrFolderInput from '../filePathInput/SubmittableFileOrFolderInput'
import { isResultObjectOk } from '../../../../preload/modules/resultStatus'

export default function imagesToDateRangeFolder(props) {
  let inputFilePathObjects = []
  let outputFilePath = ''
  const [getOutput, setGetOutput] = createSignal(function () {})
  const [status, setStatus] = createSignal('')

  async function setState(filePathObjects, path) {
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

  // We could extract similar code of the functions handleInputFilePathsRO and handleOutputDirectoryRO, for example, to the function handleRO.
  // With this extraction, handleInputFilePathsRO and handleOutputDirectoryRO call both handleRO.
  // However, this extraction hurts the performance and results in more code.
  function handleInputFilePathsRO(resultObject) {
    if (isResultObjectOk(resultObject)) {
      inputFilePathObjects = resultObject.result
    } else {
      setStatus(resultObject.message)
    }
  }

  function handleOutputDirectoryRO(resultObject) {
    if (isResultObjectOk(resultObject)) {
      // TODO: should not return an array when maxOneInput?
      outputFilePath = resultObject.result[0].value
    } else {
      setStatus(resultObject.message)
    }
  }

  function submit() {
    setGetOutput(setState(inputFilePathObjects, outputFilePath))
  }

  // TODO: minimumFiles should be 0 so it can only sort the files in destination path?
  // TODO: minimumFiles is useless in FileOrFolderInput?
  // TODO: submit should not always be part of FileOrFolderInput
  // TODO: should not be ActiveByNumberButton
  const inputComponent = (
    <div>
      <FileOrFolderInput onChange={handleInputFilePathsRO} />
      <FileOrFolderInput
        onChange={handleOutputDirectoryRO}
        filePathSelectionType={filePathSelectionType.directory}
        maxOneInput
      />
      <ActiveByNumberButton minimumNumber={1} currentNumber={1} onAction={submit} text="submit" />
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
