import FileOrFolderInput from './FileOrFolderInput'
import Page from './Page'

export default function ResultByFilesPage(props) {
  function handleFilePaths(filePaths) {
    props.onLoading(true)
    props.handleFilePaths(filePaths) // TODO: Maybe await is needed, which only works in an async function
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
