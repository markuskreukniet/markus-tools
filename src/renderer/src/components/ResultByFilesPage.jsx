import FileOrFolderInput from '../FileOrFolderInput'
import Page from './Page'

export default function ResultByFilesPage(props) {
  async function handleFilePaths(filePaths) {
    props.onLoading(true)
    props.handleFilePaths(filePaths)
    props.onLoading(false)
  }

  return (
    <Page title={props.title}>
      <FileOrFolderInput onChange={handleFilePaths} />
      <h2>Result:</h2>
      {props.resultComponent}
    </Page>
  )
}
