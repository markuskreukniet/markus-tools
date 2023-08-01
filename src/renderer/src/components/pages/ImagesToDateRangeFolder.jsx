import { createSignal } from 'solid-js'
import TextResultPage from '../page/TextResultPage'
import FileOrFolderInput from '../FileOrFolderInput'
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

  function handleFilePaths(filePaths) {
    setGetOutput(setState(filePaths, outputFilePath))
  }

  // TODO: placeholder or label? TODO: minimumFiles should be 0 so it can only sort the files in destination path?
  // TODO: minimumFiles is useless in FileOrFolderInput?
  // TODO: outputFilePath
  // TODO: submit should not always be part of FileOrFolderInput
  const inputComponent = (
    <div>
      <FileOrFolderInput onChange={handleFilePaths} minimumFiles={1} />
      <input type="text" placeholder="" onChange={(e) => (outputFilePath = e.target.value)} />
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
