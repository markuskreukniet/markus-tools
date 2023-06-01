import { createSignal } from 'solid-js'
import TextResultByFilesPage from '../page/TextResultByFilesPage'

export default function ImagesToDateRangeFolder(props) {
  const [status, setStatus] = createSignal('select and submit a folder')

  async function setState(filePaths) {
    console.log(filePaths)
    setStatus('')
  }

  return (
    <TextResultByFilesPage
      title={props.title}
      output={`Status: ${status()}`}
      minimumFiles={2}
      handleFilePaths={setState}
      onLoading={props.onLoading}
    />
  )
}
