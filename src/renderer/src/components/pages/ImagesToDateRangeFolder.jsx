import { createSignal } from 'solid-js'
import FileOrFolderInput from '../FileOrFolderInput'
import ResultPage from '../page/ResultPage'

export default function ImagesToDateRangeFolder(props) {
  const [status, setStatus] = createSignal('select and submit a folder')

  async function handleFilePaths(filePaths) {
    const status = await window.dateRangeFolder.moveImagesToDateRangeFolder(filePaths)
    setStatus(status)
  }

  // TODO: placeholder or label?
  const inputComponent = (
    <div>
      <FileOrFolderInput onChange={handleFilePaths} minimumFiles={2} />
      <input type="text" placeholder="" />
    </div>
  )
  const outputComponent = <p>{status()}</p>

  return (
    <ResultPage
      title={props.title}
      inputComponent={inputComponent}
      outputComponent={outputComponent}
      getOutput={handleFilePaths}
      onLoading={props.onLoading}
    />
  )
}
