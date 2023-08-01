import { createSignal } from 'solid-js'
import TextResultPage from '../page/TextResultPage'
import FileOrFolderInput from '../FileOrFolderInput'
import { resultStatus } from '../../../../preload/modules/resultStatus'

export default function imagesToDateRangeFolder(props) {
  let outputFilePath = ''
  const [getOutput, setGetOutput] = createSignal(function () {})
  const [status, setStatus] = createSignal('')

  async function setState(filePaths, path) {
    // TODO: error handling
    const result = await window.dateRangeFolder.imagesToDateRangeFolderBE(filePaths, path)
    const status = result.status === resultStatus.ok ? 'done' : 'not done'
    setStatus(status)
  }

  function handleFilePaths(filePaths) {
    setGetOutput(setState(filePaths, outputFilePath))
  }

  // TODO: placeholder or label? TODO: minimumFiles should be 0 so it can only sort the files in destination path?
  // TODO: minimumFiles is useless in FileOrFolderInput?
  // TODO: outputFilePath
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
