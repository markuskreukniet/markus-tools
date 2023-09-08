import { createSignal } from 'solid-js'
import TextResultPage from '../page/TextResultPage'
import ActiveByNumberButton from '../ActiveByNumberButton'
import SubmittableFileOrFolderInput from '../filePathInput/SubmittableFileOrFolderInput'
import FilePathSelector from '../filePathInput/FilePathSelector'
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

  // TODO: handleInputFilePathsRO and handleOutputDirectoryRO are almost the same
  function handleInputFilePathsRO(resultObject) {
    if (isResultObjectOk(resultObject)) {
      inputFilePathObjects = resultObject.result
    } else {
      setStatus(resultObject.message)
    }
  }

  function handleOutputDirectoryRO(resultObject) {
    if (isResultObjectOk(resultObject)) {
      outputFilePath = resultObject.result.value
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
  // TODO: should not be ActiveByNumberButton and not be SubmittableFileOrFolderInput
  const inputComponent = (
    <div>
      <SubmittableFileOrFolderInput onChange={handleInputFilePathsRO} minimumFiles={1} />
      <FilePathSelector onChange={handleOutputDirectoryRO} directory />
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
