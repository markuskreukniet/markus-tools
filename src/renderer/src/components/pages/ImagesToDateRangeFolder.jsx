import { createSignal } from 'solid-js'
import TextResultPage from '../page/TextResultPage'
import FileOrFolderInput from '../FileOrFolderInput'

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

  return (
    <TextResultPage
      title={props.title}
      inputComponent={inputComponent}
      output={status()}
      getOutput={handleFilePaths}
      onLoading={props.onLoading}
    />
  )
}
