import { createSignal } from 'solid-js'
import TextResultPage from '../page/TextResultPage'
import FileOrFolderInput from '../FileOrFolderInput'

export default function ImagesToDateRangeFolder(props) {
  let resultPath = ''
  const [getOutput, setGetOutput] = createSignal(function () {})
  const [status, setStatus] = createSignal('')

  async function setState(filePaths, path) {
    const done = await window.dateRangeFolder.moveImagesToDateRangeFolder(filePaths, path)
    const status = done ? 'done' : 'not done'
    setStatus(status)
  }

  function handleFilePaths(filePaths) {
    setGetOutput(setState(filePaths, resultPath))
  }

  // TODO: placeholder or label?
  const inputComponent = (
    <div>
      <FileOrFolderInput onChange={handleFilePaths} minimumFiles={2} />
      <input type="text" placeholder="" onChange={(e) => (resultPath = e.target.value)} />
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
