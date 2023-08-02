import { createSignal } from 'solid-js'
import TextResultPage from '../page/TextResultPage'
import FileOrFolderInput from '../FileOrFolderInput'
import FilePathSelector from '../FilePathSelector'
import { isResultObjectOk } from '../../../../preload/modules/resultStatus'

export default function imagesToDateRangeFolder(props) {
  let outputFilePath = ''
  const [getOutput, setGetOutput] = createSignal(function () {})
  const [status, setStatus] = createSignal('')

  async function setState(filePaths, path) {
    const imagesToDateRangeFolderRO = await window.dateRangeFolder.imagesToDateRangeFolderBE(
      filePaths,
      path
    )
    setStatus(
      isResultObjectOk(imagesToDateRangeFolderRO) ? 'done' : imagesToDateRangeFolderRO.message
    )
  }

  function handleInputFilePathsRO(resultObject) {
    if (isResultObjectOk(resultObject)) {
      setGetOutput(setState(resultObject.result, outputFilePath))
    } else {
      // TODO
    }
  }

  function handleOutputDirectoryRO(resultObject) {
    if (isResultObjectOk(resultObject)) {
      outputFilePath = resultObject.result
    } else {
      // TODO
    }
  }

  // TODO: minimumFiles should be 0 so it can only sort the files in destination path?
  // TODO: minimumFiles is useless in FileOrFolderInput?
  // TODO: submit should not always be part of FileOrFolderInput
  // TODO: should not select a file, but a combination a filepath and filetype (folder or file), which is possible since we have a select folder and select file button
  const inputComponent = (
    <div>
      <FileOrFolderInput onChange={handleInputFilePathsRO} minimumFiles={1} />
      <FilePathSelector onChange={handleOutputDirectoryRO} directory />
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
