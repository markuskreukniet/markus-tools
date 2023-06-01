import { createSignal } from 'solid-js'
import TextResultByFilesPage from '../page/TextResultByFilesPage'

export default function ImagesToDateRangeFolder(props) {
  const [status, setStatus] = createSignal('select and submit a folder')

  async function setState(filePaths) {
    const status = await window.dateRangeFolder.moveImagesToDateRangeFolder(filePaths)
    setStatus(status)
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
