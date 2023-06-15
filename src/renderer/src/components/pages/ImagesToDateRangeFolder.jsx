import { createSignal } from 'solid-js'
import TextResultPage from '../page/TextResultPage'
import FileOrFolderInput from '../FileOrFolderInput'

export default function ImagesToDateRangeFolder(props) {
  const [getOutput, setGetOutput] = createSignal(function () {})
  const [status, setStatus] = createSignal('')

  async function setState(filePaths) {
    const status = await window.dateRangeFolder.moveImagesToDateRangeFolder(filePaths)
    setStatus(status)
  }

  function handleFilePaths(filePaths) {
    setGetOutput(setState(filePaths))
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
      getOutput={getOutput}
      onLoading={props.onLoading}
    />
  )
}
